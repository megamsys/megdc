#!/bin/bash


########### IN Slave NODE #################

ceph_user="cibadmin"

#Install cephin slave systems
echo deb http://ceph.com/debian-giant/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
sudo apt-get -y update
sudo apt-get -y install ceph-deploy ceph-common ceph-mds

sudo apt-get -y install libvirt-bin

#Make osd directory for osd in slave systems
for d in /storage*/ ; do
	sudo mkdir $d/osd
done

#Make .ssh directory for ceph user, later master will sshpass the pub key
[ -d /home/$ceph_user/.ssh ] || mkdir /home/$ceph_user/.ssh

[ -d /home/$ceph_user/ceph-cluster ] || mkdir /home/$ceph_user/ceph-cluster

