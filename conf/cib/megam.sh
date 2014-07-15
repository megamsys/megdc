#!/bin/bash

echo "--------------------megam install----------------"

if [ $? -ne 0 ] # Did the command work?
then # Fail
     echo "An error occured. Error: \n"
     exit 1
fi

apt-get -y install opennebula >> /var/log/opennebula.log

apt-get -y remove --auto-remove opennebula >> /var/log/opennebula.log

apt-get -y purge --auto-remove opennebula  >> /var/log/opennebula.log
