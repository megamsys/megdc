#!/bin/bash
ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

echo "one_preinstall.sh start execution ====>" >> $ONE_INSTALL_LOG

ping -c 1 us.archive.ubuntu.com &> /dev/null

if [ $? -ne 0 ]; then
  echo "`date`: check your network connection. us.archive.ubuntu.com is not reachable!" >> $ONE_INSTALL_LOG
  exit 1
fi

sudo apt-get -y install build-essential autoconf libtool make >> $ONE_INSTALL_LOG

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured during the install of buildessentials.\n"
     exit 2
fi

sudo apt-get -y install lvm2 ssh iproute iputils-arping >> $ONE_INSTALL_LOG

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured during the install of ssh and bridge-utils.\n"
     exit 2
fi
echo "one_preinstall.sh end execution ====>" >> $ONE_INSTALL_LOG
