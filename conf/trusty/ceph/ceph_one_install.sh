#!/bin/bash

#bash ceph_one.sh

ceph_user="cibadmin"
host=`hostname`
poolname="one"

CEPH_INSTALL_LOG="/var/log/megam/megamcib/ceph.log"

echo "Setting up datastore for ceph in opennebula" >> $CEPH_INSTALL_LOG


sudo -H -u $ceph_user bash -c "ceph osd pool create $poolname 256"

cd /tmp

sudo -H -u $ceph_user bash -c "cat > $user_home/ds.conf <<EOF
NAME = \"cephds\"
DS_MAD = ceph
TM_MAD = ceph
DISK_TYPE = RBD
POOL_NAME = $poolname
BRIDGE_LIST = $host
CEPH_HOST = $host
EOF"

onedatastore create $user_home/ds.conf

ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=$poolname'
ceph auth get-key client.libvirt | tee client.libvirt.key
ceph auth get client.libvirt -o ceph.client.libvirt.keyring


sudo cp ceph.client.* /etc/ceph

uid=`uuidgen`

sudo cat > secret.xml <<EOF
<secret ephemeral='no' private='no'>
  <uuid>$uid</uuid>
  <usage type='ceph'>
          <name>client.libvirt secret</name>
  </usage>
</secret>
EOF

sudo apt-get -y install libvirt-bin

sudo virsh secret-define secret.xml

sudo virsh secret-set-value --secret $uid --base64 $(cat client.libvirt.key)

#Update datastore for ceph
#onedatastore show cephds | grep "ID "
##UPDATE IT IN ONE_UI
#CEPH_USER="libvirt"
#CEPH_SECRET="$UUID"
#CEPH_HOST="<list of ceph mon hosts, see table above>"






