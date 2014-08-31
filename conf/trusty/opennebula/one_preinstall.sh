#!/bin/bash

ping -c 1 us.archive.ubuntu.com &> /dev/null

if [ $? -ne 0 ]; then
  echo "`date`: check your network connection. us.archive.ubuntu.com is not reachable!" >> /var/log/megam/megamcib/one_preinstall.log
  exit 1
fi

sudo apt-get -y install qemu-system-x86 build-essential genromfs autoconf libtool qemu-utils libvirt0 bridge-utils qemu-kvm

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured during the install of qemu, buildessentials.\n"
     exit 2
fi

sudo apt-get -y install lvm2 ssh bridge-utils iproute iputils-arping make

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured during the install of ssh and bridge-utils.\n"
     exit 2
fi
