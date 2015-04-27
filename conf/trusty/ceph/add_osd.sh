#!/bin/bash


########### IN Slave NODE #################

ceph_user="cibadmin"
CEPH_LOG="/var/log/megam/megamcib/ceph.log"

#Install cephin slave systems
echo "add_osd.sh start execution ====>" >> $CEPH_LOG
echo deb http://ceph.com/debian-hammer/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
sudo apt-get -y update
sudo apt-get -y install ceph-deploy ceph-common ceph-mds >> $CEPH_LOG

sudo apt-get -y install libvirt-bin >> $CEPH_LOG

#Make osd directory for osd in slave systems
for d in /storage*/ ; do
	sudo mkdir $d/osd
	echo "Created Directory $d/osd ====>" >> $CEPH_LOG
done

#Make .ssh directory for ceph user, later master will sshpass the pub key
[ -d /home/$ceph_user/.ssh ] || mkdir /home/$ceph_user/.ssh

[ -d /home/$ceph_user/ceph-cluster ] || mkdir /home/$ceph_user/ceph-cluster

echo "add_osd.sh end execution ====>" >> $CEPH_LOG
