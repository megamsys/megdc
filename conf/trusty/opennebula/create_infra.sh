#!/bin/bash
###############################
#### HOST CREATION ############
###############################
sudo -H -u oneadmin bash -c 'onehost create $1 -i kvm -v kvm -n dummy'

ipaddr = $1

#Get first 3 values of ip4 eg:192.168.1 in 192.168.1.100
ip3=`echo $ipaddr| cut -d'.' -f 1,2,3`

###############################
#### NETWORK CREATION #########
###############################

cat > //var/lib/one/vn.net <<EOF
NAME   = "open-vs"
TYPE   = FIXED
BRIDGE = one
AR = [ TYPE = "IP4", IP   = "$ip3.120", SIZE = "100" ]
DNS = "8.8.8.8 8.8.4.4"
GATEWAY    = "$ip3.1"
EOF

sudo chown oneamdin:oneamdin /var/lib/one/vn.net

sudo -H -u oneadmin bash -c 'onevnet create /var/lib/one/vn.net'



