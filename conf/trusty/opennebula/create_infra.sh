#!/bin/bash

#bash create_infra.sh UID
host=`hostname`
#### HOST CREATION ############
sudo -H -u oneadmin bash -c 'onehost create $host -i kvm -v kvm -n dummy'


#### DATASTORE CREATION ############
sudo -H -u oneadmin bash -c "cat > /var/lib/one/ds.conf <<EOF
NAME = \"cephds\"
DS_MAD = ceph
TM_MAD = ceph
DISK_TYPE = RBD
CEPH_USER = libvirt
CEPH_SECRET = $uid
POOL_NAME = $poolname
BRIDGE_LIST = $host
CEPH_HOST = $host
EOF"

sudo -H -u oneadmin bash -c 'onedatastore create /var/lib/one/ds.conf'
echo "Setting up datastore for ceph in opennebula" >> $CEPH_INSTALL_LOG


#### NETWORK CREATION #########

while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
done < /proc/net/route


#Get first 3 values of ip4 eg:192.168.1 in 192.168.1.100
ip3=`echo $ipaddr| cut -d'.' -f 1,2,3`


cat > //var/lib/one/vn.net <<EOF
NAME   = "open-vs"
TYPE   = FIXED
BRIDGE = one
AR = [ TYPE = "IP4", IP   = "$ip3.206", SIZE = "50" ]
DNS = "8.8.8.8 8.8.4.4"
GATEWAY    = "$ip3.1"
EOF

sudo chown oneamdin:oneamdin /var/lib/one/vn.net

sudo -H -u oneadmin bash -c 'onevnet create /var/lib/one/vn.net'



