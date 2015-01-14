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
listen     megamnialvu         192.168.1.100:8080
                 mode http
                 stats enable
                 stats auth cibadmin:cibadmin # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server megammaster 192.168.1.100:8080 cookie A check
                 server megamslave 192.168.1.101:8080 cookie B check


listen     opennebula         192.168.1.100:9869
                 mode http
                 stats enable
                 stats auth cibadmin:cibadmin # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server megammaster 192.168.1.100:9869 cookie A check
                 server megamslave 192.168.1.101:9869 cookie B check

listen     apache         192.168.1.100:80
                 mode http
                 stats enable
                 stats auth cibadmin:cibadmin # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server megammaster 192.168.1.100:80 cookie A check
                 server megamslave 192.168.1.101:80 cookie B check
EOT

service haproxy restart

