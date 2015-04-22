#!/bin/bash

MEGAM_LOG="/var/log/megam/megamcib/megam.log"

ping -c 1 get.megam.io &> /dev/null

if [ $? -ne 0 ]; then
	echo "`date`: check your network connection. get.megam.io is down or not reachable!" >> $MEGAM_LOG
  exit 1
fi


echo "Adding entries in /etc/hosts" >> $MEGAM_LOG


	while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route

echo "127.0.0.1 `hostname` localhost" >> /etc/hosts
echo "$ipaddr `hostname` localhost" >> /etc/hosts


apt-get -y install megamcommon >> $MEGAM_LOG

apt-get -y install megamnilavu >> $MEGAM_LOG

sudo apt-get -y install software-properties-common python-software-properties >> $MEGAM_LOG

sudo apt-add-repository -y ppa:openjdk-r/ppa >> $MEGAM_LOG

sudo apt-get -y update >> $MEGAM_LOG

sudo apt-get -y install openjdk-8-jdk >> $MEGAM_LOG

apt-get -y install megamgateway >> $MEGAM_LOG

apt-get -y install chef-server >> $MEGAM_LOG

apt-get -y install rabbitmq-server >> $MEGAM_LOG

apt-get -y install megamd >> $MEGAM_LOG

apt-get -y install megamanalytics >> $MEGAM_LOG

apt-get -y install megamchefnative >> $MEGAM_LOG

echo "`date`: Step1: megam installed successfully." >> $MEGAM_LOG
