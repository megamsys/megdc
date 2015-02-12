#!/bin/bash

#bash ceph_install.sh install osd1="/storage1" osd2="/storage2" osd3="/storage3"

for i in "$@"
do
case $i in
    osd1=*)
    osd1="${i#*=}"
    ;;
    osd2=*)
    osd2="${i#*=}"
    ;;
    osd3=*)
    osd3="${i#*=}"
    ;;   
esac
done

echo "=================ceph storages======================="
echo $osd1
echo $osd2
echo $osd3
echo "========================================"


MEGAM_LOG="/var/log/megam/megamcib/megam.log"
echo "Step 1: you are running the test file. megam install is sleeping..." >> $MEGAM_LOG
sleep 10
exit 0
