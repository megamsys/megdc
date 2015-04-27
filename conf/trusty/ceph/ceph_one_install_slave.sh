#!/bin/bash

poolname="one"
user_home="/home/cibadmin"

CEPH_LOG="/var/log/megam/megamcib/ceph.log"

echo "ceph_one_install_slave.sh start execution ====>" >> $CEPH_LOG
ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=$poolname'

cd $user_home/ceph-cluster

echo "virsh secret-define secret.xml" >> $CEPH_LOG
sudo virsh secret-define secret.xml

uid=`cat uid`

echo "virsh secret-set value " >> $CEPH_LOG
sudo virsh secret-set-value --secret $uid --base64 $(cat client.libvirt.key)


#After, this executiopn ends, Add host in opennebula master
#And also change datastore 'cephds'(Add slave host in 'CEPH_HOST')

echo "ceph_one_install_slave.sh end execution ====>" >> $CEPH_LOG
