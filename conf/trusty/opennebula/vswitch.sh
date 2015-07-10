#!/bin/bash

sudo apt-get -y install openvswitch-common openvswitch-switch bridge-utils


cat > //etc/sudoers.d/openvswitch <<EOF
%oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-vsctl
%oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-ofctl
EOF

BRIDGE_NAME="one"
NETWORK_IF="eth0"

#sudo ovs-vsctl add-br $BRIDGE_NAME

echo "auto $BRIDGE_NAME" >> /etc/network/interfaces
#sudo ovs-vsctl show

#ovs-vsctl add-port $BRIDGE_NAME $NETWORK_IF

#sudo ifconfig one up

#sudo ovs-vsctl add-port one eth0
#sudo ovs-vsctl del-port one eth0


#dhclient one

#=================================================================================================
#change broadcast 'eth0' as 'one' in /etc/ha.d/ha.cf for drbd heartbeat
#=================================================================================================
ip()
{
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

#Get first 3 values of ip4 eg:192.168.1 in 192.168.1.100
ip3=`echo $ipaddr| cut -d'.' -f 1,2,3`


cat > //etc/network/interfaces <<EOF
auto eth0
#iface eth0 inet static

auto one
iface one inet static
address $ip3.100
network $ip3.0
gateway $ip3.1
netmask 255.255.255.0
broadcast $ip3.255
bridge_ports eth0
dns-nameservers 8.8.8.8 8.8.4.4
EOF

#ifconfig eth0 0			#=====> Connection will be cleared

#System needs to be restarted



