#!/bin/bash

#haproxy.sh node1=node1:port node2=node2:port

ip() {
while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route
}

ip

for i in "$@"
do
case $i in
    node1=*)
    node1="${i#*=}"
    ;;
    node2=*)
    node2="${i#*=}"
    ;;
esac
done


gateway="$ipaddr:80"

address1="$node1"
address2="$node2"

sudo add-apt-repository -y ppa:vbernat/haproxy-1.5
sudo apt-get -y update || true
apt-get install -y haproxy

echo "net.ipv4.ip_nonlocal_bind=1" >> /etc/sysctl.conf
echo "ENABLED=1" >> /etc/default/haproxy

mv /etc/haproxy/haproxy.cfg{,.original}

cat << EOT >> /etc/haproxy/haproxy.cfg
global
	log /dev/log	local0
	log /dev/log	local1 notice
	chroot /var/lib/haproxy
	stats socket /run/haproxy/admin.sock mode 660 level admin
	stats timeout 30s
	user haproxy
	group haproxy
	daemon

	# Default SSL material locations
	ca-base /etc/ssl/certs
	crt-base /etc/ssl/private

	# Default ciphers to use on SSL-enabled listening sockets.
	# For more information, see ciphers(1SSL).
	ssl-default-bind-ciphers kEECDH+aRSA+AES:kRSA+AES:+AES256:RC4-SHA:!kEDH:!LOW:!EXP:!MD5:!aNULL:!eNULL
        ssl-default-bind-options no-sslv3

defaults
	log	global
	mode	http
	option	httplog
	option	dontlognull
        timeout connect 5000
        timeout client  50000
        timeout server  50000
	errorfile 400 /etc/haproxy/errors/400.http
	errorfile 403 /etc/haproxy/errors/403.http
	errorfile 408 /etc/haproxy/errors/408.http
	errorfile 500 /etc/haproxy/errors/500.http
	errorfile 502 /etc/haproxy/errors/502.http
	errorfile 503 /etc/haproxy/errors/503.http
	errorfile 504 /etc/haproxy/errors/504.http


listen     web-cluster         $gateway
                 mode http
                 stats enable
                 stats auth cibadmin:cibadmin # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server megammaster $address1 cookie A check
                 server megamslave $address2 cookie B check
EOT

service haproxy restart

