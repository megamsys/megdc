#!/bin/bash

wget -q -O- http://downloads.opennebula.org/repo/Ubuntu/repo.key | apt-key add -
echo "deb http://downloads.opennebula.org/repo/Ubuntu/14.04 stable opennebula" > /etc/apt/sources.list.d/opennebula.list

apt-get update

apt-get upgrade

sudo apt-get install opennebula opennebula-sunstone
