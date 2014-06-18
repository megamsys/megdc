#!/bin/bash

sudo apt-get -y install qemu-system-x86 build-essential genromfs autoconf libtool qemu-utils libvirt0 bridge-utils qemu-kvm

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

sudo apt-get -y install lvm2 ssh bridge-utils iproute iputils-arping make m4 \
                  ndisc6 python python-openssl openssl \
                  python-pyparsing python-simplejson python-bitarray \
                  python-pyinotify python-pycurl python-ipaddr socat fping

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

sudo apt-get -y install python-paramiko python-affinity qemu-utils

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

#Python tools and libraries
sudo apt-get -y install python-setuptools python-dev automake git fakeroot python-yaml pylint pep8

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

cd / && easy_install affinity bitarray ipaddr

#Haskell Install
sudo apt-get -y install ghc libghc-json-dev libghc-network-dev \
                  libghc-parallel-dev libghc-deepseq-dev \
                  libghc-utf8-string-dev libghc-curl-dev \
                  libghc-hslogger-dev \
                  libghc-crypto-dev libghc-text-dev \
                  libghc-hinotify-dev libghc-regex-pcre-dev \
                  libpcre3-dev \
                  libghc-attoparsec-dev libghc-vector-dev \
		  libghc-zlib-dev libghc-base64-bytestring-dev

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

sudo apt-get -y install libghc-quickcheck2-dev libghc-hunit-dev \
      libghc-test-framework-dev \
      libghc-test-framework-quickcheck2-dev \
      libghc-test-framework-hunit-dev \
      libghc-temporary-dev shelltestrunner \
      hscolour hlint

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

#DRBD installation
sudo apt-get -y install drbd8-utils

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

sudo -s && echo drbd  minor_count=128 usermode_helper=/bin/true >> /etc/modules

depmod -a

modprobe drbd minor_count usermode_helper=/bin/true





