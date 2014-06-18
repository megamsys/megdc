#!/bin/bash

#input for node and cluster
CLUSTERNAME="$0"
NODEIP="$1"
NODENAME="$2"

#setup key pair for cluster
sudo ssh-keygen

sudo mkdir -m 700 -p /root/.ssh && ln -s /etc/ssh/ssh_host_rsa_key /root/.ssh/id_rsa

brctl addbr br0
brctl addif br0 eth0

#Initialize Cluster
sudo gnt-cluster init --master-netdev br0 --vg-name xenvg --enabled-hypervisors kvm --nic-parameters link=br0 --mac-prefix 00:16:37 --no-ssh-init --no-etc-hosts --hypervisor-parameters kvm:initrd_path=,kernel_path= "${CLUSTERNAME}" 

#Verify Cluster initialization
sudo gnt-cluster verify 

#List nodes
sudo gnt-node list

#Peparation for boostrapping images.
wget http://ganeti.googlecode.com/files/ganeti-instance-debootstrap-0.14.tar.gz
tar xzf ganeti-*.tar.gz
cd ganeti-instance-debootstrap-0.14
./configure --with-os-dir=/srv/ganeti/os

make
sudo make install

sudo apt-get install debootstrap dump kpartx

sudo gnt-os list

#Add trusty image
echo "trusty" >> /srv/ganeti/os/debootstrap/variants.list

cat > /srv/ganeti/os/debootstrap/variants/trusty.conf <<EOF
MIRROR="http://archive.ubuntu.com/ubuntu/"
SUITE="trusty"
COMPONENTS="main,universe"
ARCH="amd64"
EXTRA_PKGS="acpi-support,udev,linux-image-generic-lts-trusty,grub2,openssh-server,curl"  
EOF

cat > /usr/local/etc/ganeti/instance-debootstrap/variants/trusty.conf <<EOF
MIRROR="http://archive.ubuntu.com/ubuntu/"
SUITE="trusty"
COMPONENTS="main,universe"
ARCH="amd64"
EXTRA_PKGS="acpi-support,udev,linux-image-generic-lts-trusty,grub2,openssh-server,curl"  
EOF

sudo cp -r $(dirname $0)/hooks/* /usr/local/etc/ganeti/instance-debootstrap/hooks/

sudo find /usr/local/etc/ganeti/instance-debootstrap/hooks/ -type f -exec chmod 754 {} \;

sudo gnt-node add -g default -s "${NODEIP}" "${NODENAME}"








