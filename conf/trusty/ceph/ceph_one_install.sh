#!/bin/bash

black='\033[30m'
red='\033[31m'
green='\033[32m'
yellow='\033[33m'
blue='\033[34m'
magenta='\033[35m'
cyan='\033[36m'
white='\033[37m'

alias Reset="tput sgr0"      #  Reset text attributes to normal
# without clearing screen.


#--------------------------------------------------------------------------
#colored echo
# Argument $1 = message
# Argument $2 = color (
#--------------------------------------------------------------------------
cecho () {
  local default_msg="No message passed."  # Doesn't really need to be a local variable.
  message=${1:-$default_msg}              # Defaults to default message.
  color=${2:-$black}                      # Defaults to black, if not specified.
  echo "$color$message"
  Reset                                   # Reset to normal.
  return
}




#bash ceph_one.sh

ceph_user="cibadmin"
host=`hostname`
poolname="one"

user_home="/home/cibadmin"

CEPH_INSTALL_LOG="/var/log/megam/megamcib/ceph.log"


echo "Creating ceph osd pool... $poolname " >> $CEPH_INSTALL_LOG
sudo -H -u $ceph_user bash -c "ceph osd pool create $poolname 256"



sudo -H -u oneadmin bash -c "cat > /var/lib/one/ds.conf <<EOF
NAME = \"cephds\"
DS_MAD = ceph
TM_MAD = ceph
DISK_TYPE = RBD
POOL_NAME = $poolname
BRIDGE_LIST = $host
CEPH_HOST = $host
EOF"

sudo -H -u oneadmin bash -c 'onedatastore create /var/lib/one/ds.conf'
echo "Setting up datastore for ceph in opennebula" >> $CEPH_INSTALL_LOG

cd /tmp

echo "processing get-or-create auth user..." >> $CEPH_INSTALL_LOG
ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=$poolname'
ceph auth get-key client.libvirt | tee client.libvirt.key
ceph auth get client.libvirt -o ceph.client.libvirt.keyring

echo "Copying keyrings to /etc/ceph..." >> $CEPH_INSTALL_LOG
sudo cp ceph.client.* /etc/ceph

uid=`uuidgen`

echo "creating secret.xml file in /tmp dir..." >> $CEPH_INSTALL_LOG

sudo cat > secret.xml <<EOF
<secret ephemeral='no' private='no'>
  <uuid>$uid</uuid>
  <usage type='ceph'>
          <name>client.libvirt secret</name>
  </usage>
</secret>
EOF


sudo apt-get -y install libvirt-bin >> $CEPH_INSTALL_LOG

echo "virsh secret-define secret.xml" >> $CEPH_INSTALL_LOG
sudo virsh secret-define secret.xml >> $CEPH_INSTALL_LOG

echo "virsh secret-define secret.xml" >> $CEPH_INSTALL_LOG
sudo virsh secret-set-value --secret $uid --base64 $(cat client.libvirt.key)

#Update datastore for ceph
#onedatastore show cephds | grep "ID "
##UPDATE IT IN ONE_UI
#CEPH_USER="libvirt"
#CEPH_SECRET="$UUID"
#CEPH_HOST="<list of ceph mon hosts"







