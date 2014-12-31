#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

echo "oneadmin ALL = (root) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/oneadmin            #all-nodes

sudo apt-get -y install ntp                                                                  #all-nodes

sudo apt-get -y install ruby-dev

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
#Add ip and port of sunstone-server in conf
sed -i "s/^[ \t]*:host:.*/:host: $ipaddr/" /etc/one/sunstone-server.conf

#Cahnge all datastore's tm_mad to ssh
sudo -H -u oneadmin bash -c "echo \"TM_MAD=ssh\" > /tmp/ds_tm_mad"
sudo -H -u oneadmin bash -c "onedatastore update 0 /tmp/ds_tm_mad"
sudo -H -u oneadmin bash -c "onedatastore update 1 /tmp/ds_tm_mad"
sudo -H -u oneadmin bash -c "onedatastore update 2 /tmp/ds_tm_mad"

#Edit clone file for scp problem
sed -i '/SRC=$1/a SRC=${SRC#*:}' /var/lib/one/remotes/tm/ssh/clone

service_restart() {
sunstone-server restart
econe-server restart
occi-server restart
onegate-server restart
one restart
}

service_restart





#Add code to automate template, image, and network



