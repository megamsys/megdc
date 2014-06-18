#!/bin/bash

sudo apt-get -y install qemu-system-x86 build-essential genromfs autoconf libtool qemu-utils libvirt0 bridge-utils qemu-kvm

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

sudo apt-get -y install lvm2 ssh bridge-utils iproute iputils-arping make

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi
