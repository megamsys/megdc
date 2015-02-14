#!/bin/bash

#haproxy.sh node1_ip=node1_ip node2_ip=node2_ip node1_host=node1_host node2_host=node2_host

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
    node1_ip=*)
    node1_ip="${i#*=}"
    ;;
    node2_ip=*)
    node2_ip="${i#*=}"
    ;;
    node1_host=*)
    node1_host="${i#*=}"
    ;;
    node2_host=*)
    node2_host="${i#*=}"
    ;;
esac
done

HA_LOG="/var/log/megam/megamcib/ha.log"
gateway="$ipaddr"

address1="$node1_ip"
address2="$node2_ip"

host1="$node1_host"
host2="$node2_host"

username="cibadmin"
password="cibadmin"
auth_credential=$username:$password

sudo add-apt-repository -y ppa:vbernat/haproxy-1.5
sudo apt-get -y update || true
apt-get install -y haproxy >> $HA_LOG

echo "net.ipv4.ip_nonlocal_bind=1" >> /etc/sysctl.conf
echo "ENABLED=1" >> /etc/default/haproxy

mv /etc/haproxy/haproxy.cfg{,.original}

echo "Writing /etc/haproxy/haproxy.conf"  >> $HA_LOG
cat << EOT >> /etc/haproxy/haproxy.cfg
listen     megamnialvu         $gateway:8080
                 mode http
                 stats enable
                 stats auth $auth_credential # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server $host1 $address1:8080 cookie A check
                 server $host2 $address2:8080 cookie B check


listen     opennebula         $gateway:9869
                 mode http
                 stats enable
                 stats auth $auth_credential # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server $host1 $address1:9869 cookie A check
                 server $host2 $address2:9869 cookie B check

listen     apache         $gateway:80
                 mode http
                 stats enable
                 stats auth $auth_credential # Change this to your own username and password!
                 balance roundrobin
                 option httpclose
                 option forwardfor
                 cookie JSESSIONID prefix
                 server $host1 $address1:80 cookie A check
                 server $host2 $address2:80 cookie B check
EOT

echo "Restart haproxy service"  >> $HA_LOG
service haproxy restart

