#!/bin/bash

ping -c 1 us.archive.ubuntu.com &> /dev/null

if [ $? -ne 0 ]; then
  echo "`date`: check your network connection. us.archive.ubuntu.com is not reachable!" >> /var/log/megam/megamcib/opennebula.log
  exit 1
fi

sudo apt-get -y install build-essential autoconf libtool make

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured during the install of buildessentials.\n"
     exit 2
fi

sudo apt-get -y install lvm2 ssh iproute iputils-arping 

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured during the install of ssh and bridge-utils.\n"
     exit 2
fi
