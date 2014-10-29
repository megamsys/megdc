#!/bin/bash

echo deb http://ceph.com/debian-giant/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
sudo apt-get -y update
sudo apt-get -y install ceph-deploy ceph-common ceph-mds

echo "oneadmin ALL = (root) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/oneadmin            #all-nodes
sudo apt-get install ntp                                                                  #all-nodes
sudo chmod 0440 /etc/sudoers.d/oneadmin                                                  #all-nodes

#nano /etc/hosts

echo "192.168.2.10 mon0" >> /etc/hosts                  #admin-node ip and hostname
echo "192.168.2.8 osd0" >> /etc/hosts                   #storage-node-1 ip and hostname
echo "192.168.2.8 osd1" >> /etc/hosts                   #storage-node2 ip and hostname

sudo apt-get install dnsmasq 


#nano .ssh/config
cat > //varr/lib/one/.ssh/config <<EOF
Host osd0
   Hostname trusty-megamnode-x8664                      #storage node-1
   User oneadmin
Host osd1
   Hostname trusty-megamnode-x8664                      #storage node-2
   User oneadmin
Host mon0
   Hostname megamubuntu                                 #Admin hostname
   User oneadmin
EOF


su oneadmin

cd ~

mkdir ceph-cluster

cd ceph-cluster

ceph-deploy new megamubuntu             #[admin-node]
ceph-deploy install --no-adjust-repos megamubuntu osd0 osd1

#PROMPT  Are you sure you want to continue connecting (yes/no)? yes    fir the first time          oneadmin@megamubuntu's password: 

#ceph-deploy mon create-initial          ---> doesnot get ceph.bootstrap-mds.keyring

sudo ceph-deploy mon create-initial
sudo chown oneadmin:oneadmin *.keyring

ceph-deploy osd --fs-type ext4 prepare osd0:/ceph-store/osd osd1:/ceph-store/osd
ceph-deploy osd --fs-type ext4 activate osd0:/ceph-store/osd osd1:/ceph-store/osd

sudo chmod +r /etc/ceph/ceph.client.admin.keyring

ceph-deploy admin megamubuntu osd0 osd1
ceph-deploy mds create megamubuntu


ceph osd pool set metadata size 2
ceph osd pool set data size 2

sudo stop ceph-all
sudo start ceph-all


