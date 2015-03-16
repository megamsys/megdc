#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

echo "Installing sshpass " >> $ONE_INSTALL_LOG
sudo apt-get -y install sshpass


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

#if front-end and host are seperate servers
if [ $ipaddr != $1]
then
#sudo rm /var/lib/one/.ssh/id_rsa*		#If already existed key is used by anyother systems?
echo "Generating ssh_key for oneadmin user " >> $ONE_INSTALL_LOG
#sudo -H -u oneadmin bash -c "ssh-keygen -N '' -t rsa -f /var/lib/one/.ssh/id_rsa"
fi


echo "Transfering auth_keys to megamcib_node " >> $ONE_INSTALL_LOG
sshpass -p "oneadmin" scp -o StrictHostKeyChecking=no /var/lib/one/.ssh/id_rsa.pub oneadmin@$1:/var/lib/one/.ssh/authorized_keys
#sshpass -p "oneadmin" scp -o StrictHostKeyChecking=no /var/lib/one/.ssh/id_rsa.pub oneadmin@192.168.6.201:/var/lib/one/.ssh/authorized_keys

#No prompt on "Add ip to known_hosts list"
sudo -H -u oneadmin bash -c "cat > //var/lib/one/.ssh/ssh_config <<EOF
ConnectTimeout 5
Host *
StrictHostKeyChecking no
EOF"

echo "Oneadmin Authenticated. Oneadmin can access hosts without password "

#onehost create 192.168.6.201 -i kvm -v kvm -n dummy
#su oneadmin




