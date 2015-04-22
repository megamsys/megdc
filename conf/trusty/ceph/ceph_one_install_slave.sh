#!/bin/bash

poolname="one"
user_home="/home/cibadmin"


ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=$poolname'

cd $user_home/ceph-cluster

sudo virsh secret-define secret.xml

uid=`cat uid`

sudo virsh secret-set-value --secret $uid --base64 $(cat client.libvirt.key)


#After, this executiopn ends, Add host in opennebula master
#And also change datastore 'cephds'(Add slave host in 'CEPH_HOST')
