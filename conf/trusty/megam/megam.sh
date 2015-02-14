#!/bin/bash

MEGAM_LOG="/var/log/megam/megamcib/megam.log"

ping -c 1 get.megam.co &> /dev/null

if [ $? -ne 0 ]; then
	echo "`date`: check your network connection. get.megam.co is down or not reachable!" >> $MEGAM_LOG
  exit 1
fi


echo "Installing megam.."

apt-get -y install megamcommon >> $MEGAM_LOG

apt-get -y install megamnilavu >> $MEGAM_LOG

apt-get -y install megamgateway >> $MEGAM_LOG

apt-get -y install chef-server >> $MEGAM_LOG

apt-get -y install megamd >> $MEGAM_LOG

apt-get -y install megamanalytics >> $MEGAM_LOG

apt-get -y install megamchefnative >> $MEGAM_LOG

echo "`date`: Step1: megam installed successfully." >> $MEGAM_LOG
