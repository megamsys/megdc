#!/bin/bash

#sh megam.sh remote_ip="192.168.2.101" remote_hostname="megamslave" data_dir="/var/lib/" local_disk="/dev/sda9" remote_disk="/dev/sda9" master

#Get ip of the two nodes as argument
master=false
for i in "$@"
do
case $i in
    remote_ip=*)
    remote_ip="${i#*=}"
    ;;
    remote_hostname=*)
    remote_hostname="${i#*=}"
    ;;
    remote_disk=*)
    remote_disk="${i#*=}"
    ;;
    local_disk=*)
    local_disk="${i#*=}"
    ;;
    data_dir=*)
    data_dir="${i#*=}"
    ;;
    master*)
    master=true
    ;;
esac
done

echo "===============megam.sh========================="
echo $remote_ip 
echo $remote_hostname
echo $remote_disk
echo $local_disk
echo "========================================"

MEGAM_LOG="/var/log/megam/megamcib/megam.log"
echo "Step 1: you are running the test file. megam install is sleeping..." >> $MEGAM_LOG
sleep 10
exit 0
