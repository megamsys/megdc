#!/bin/bash

CEPH_INSTALL_LOG="/var/log/megam/megamcib/ceph.log"

echo "Transfering auth_keys to megamcib_node " >> $CEPH_INSTALL_LOG
sshpass -p "oneadmin" scp /var/lib/one/.ssh/id_rsa.pub oneadmin@$1:/var/lib/one/.ssh/authorized_keys
echo "Oneadmin Authenticated. Oneadmin can access hosts without password "


