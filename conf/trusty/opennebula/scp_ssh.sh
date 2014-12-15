#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

echo "Installing sshpass " >> $ONE_INSTALL_LOG
sudo apt-get -y install sshpass

sudo rm /var/lib/one/.ssh/id_rsa*

echo "Generating ssh_key for oneadmin user " >> $ONE_INSTALL_LOG
sudo -H -u oneadmin bash -c 'ssh-keygen -t rsa'

echo "Transfering auth_keys to megamcib_node " >> $ONE_INSTALL_LOG
sshpass -p "oneadmin" scp /var/lib/one/.ssh/id_rsa.pub oneadmin@$1:/var/lib/one/.ssh/authorized_keys
echo "Oneadmin Authenticated. Oneadmin can access hosts without password "

#onehost create 192.168.2.129 -i kvm -v kvm -n dummy
#su oneadmin
sudo -H -u oneadmin bash -c 'onehost create $1 -i kvm -v kvm -n dummy'

