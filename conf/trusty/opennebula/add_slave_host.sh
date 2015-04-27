#!/bin/bash

#bash add_slave_host.sh SLAVE_IP_ADDRESS

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

echo "add_slave_host.sh start execution ====>" >> $ONE_INSTALL_LOG

#Add slave host
sudo -H -u oneadmin bash -c "onehost create $1 -i kvm -v kvm -n dummy"
echo "Host $1 created in one ====>" >> $ONE_INSTALL_LOG

#Cahnge ceph datastore's ceph_host
sudo -H -u oneadmin bash -c "echo \"CEPH_HOST = `hostname` $1\" > /tmp/ds_ceph"

#Assume ceph datastore's id 100
echo "Adding ceph Host $1 in cephds ====>" >> $ONE_INSTALL_LOG
sudo -H -u oneadmin bash -c "onedatastore update 100 /tmp/ds_ceph"  >> $ONE_INSTALL_LOG

echo "Ceph Datastore updated with host $1 ====>" >> $ONE_INSTALL_LOG

echo "One service restart start ====>" >> $ONE_INSTALL_LOG
sudo sunstone-server restart >> $ONE_INSTALL_LOG
sudo econe-server restart >> $ONE_INSTALL_LOG
sudo occi-server restart >> $ONE_INSTALL_LOG
sudo onegate-server restart >> $ONE_INSTALL_LOG
sudo -H -u oneadmin bash -c "one restart" >> $ONE_INSTALL_LOG
sudo service opennebula restart >> $ONE_INSTALL_LOG


echo "add_slave_host.sh end execution ====>" >> $ONE_INSTALL_LOG
