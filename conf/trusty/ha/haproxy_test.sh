#!/bin/bash

#haproxy.sh node1_ip=node1_ip node2_ip=node2_ip node1_host=node1_host node2_host=node2_host

for i in "$@"
do
case $i in
    node1_ip=*)
    node1_ip="${i#*=}"
    ;;
    node2_ip=*)
    node2_ip="${i#*=}"
    ;;
    node1_host=*)
    node1_host="${i#*=}"
    ;;
    node2_host=*)
    node2_host="${i#*=}"
    ;;
esac
done

echo "=================ha proxy======================="
echo $node1_ip 
echo $node1_host
echo $node2_ip
echo $node2_host
echo "========================================"

MEGAM_LOG="/var/log/megam/megamcib/megam.log"
echo "Step 1: you are running the test file. megam install is sleeping..." >> $MEGAM_LOG
sleep 10
exit 0
