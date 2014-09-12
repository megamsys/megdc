#!/bin/bash
sudo apt-get -y remove opennebula opennebula-sunstone
sudo apt-get -y purge opennebula opennebula-sunstone
sudo apt-get -y autoremove

sudo rm /etc/apt/sources.list.d/opennebula.list

