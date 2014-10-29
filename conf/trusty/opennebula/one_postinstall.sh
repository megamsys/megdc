#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

echo "oneadmin ALL = (root) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/oneadmin            #all-nodes

sudo apt-get install ntp                                                                  #all-nodes

sudo chmod 0440 /etc/sudoers.d/oneadmin                                                  #all-nodes

echo "One POSTINSTALL install gems =======> " >> $ONE_INSTALL_LOG

#sed -i 's/.*dependencies_common => .*/:dependencies_common => [],/' /usr/share/one/install_gems
sudo rm /usr/share/one/install_gems

#/var/lib/megam/megamcib/install_gems has to be created early
sudo cp /usr/share/megam/megamcib/conf/trusty/opennebula/install_gems /usr/share/one/install_gems

sudo chmod 755 /usr/share/one/install_gems

sudo /usr/share/one/install_gems sunstone >> $ONE_INSTALL_LOG

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
sed -i "s/^[ \t]*:host:.*/:host: $ipaddr/" /etc/one/sunstone-server.conf
