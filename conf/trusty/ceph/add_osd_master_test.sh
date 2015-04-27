#!/bin/bash
MEGAM_LOG="/var/log/megam/megamcib/megam.log"
echo "Step 1: you are running the test file. add osd master is sleeping..." >> $MEGAM_LOG
echo $1
sleep 10
exit 0
