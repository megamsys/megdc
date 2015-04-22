#!/bin/bash

#bash add_slave_host.sh SLAVE_IP_ADDRESS

#Add slave host
sudo -H -u oneadmin bash -c 'onehost create $1 -i kvm -v kvm -n dummy'

#Cahnge ceph datastore's ceph_host
sudo -H -u oneadmin bash -c "echo \"CEPH_HOST = `hostname` $1\" > /tmp/ds_ceph"
sudo -H -u oneadmin bash -c "onedatastore update 100 /tmp/ds_ceph"


sudo sunstone-server restart
sudo econe-server restart
sudo occi-server restart
sudo onegate-server restart
sudo -H -u oneadmin bash -c "one restart"
sudo service opennebula restart
