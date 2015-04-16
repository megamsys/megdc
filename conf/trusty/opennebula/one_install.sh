#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

ping -c 1 downloads.opennebula.org &> /dev/null

if [ $? -ne 0 ]; then
  echo "`date`: check your network connection. downloads.opennebula.org is not reachable!" >> $ONE_INSTALL_LOG
  exit 1
fi

wget -q -O- http://downloads.opennebula.org/repo/Ubuntu/repo.key | apt-key add -

echo "deb http://downloads.opennebula.org/repo/4.12/Ubuntu/14.04 stable opennebula" > /etc/apt/sources.list.d/opennebula.list

echo "File created /etc/apt/sources.list.d/opennebula.list " >> $ONE_INSTALL_LOG

sudo apt-get -y update >> $ONE_INSTALL_LOG


sudo apt-get -y install opennebula opennebula-sunstone >> $ONE_INSTALL_LOG
