#!/bin/bash

echo "--------------------megam install----------------"

apt-get -y install opennebula >> /var/log/opennebula.log

apt-get -y remove --auto-remove opennebula >> /var/log/opennebula.log

apt-get -y purge --auto-remove opennebula  >> /var/log/opennebula.log
