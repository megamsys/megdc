#!/bin/bash

#Before this execution starts, run add_osd_slave.sh in slave system.

########### IN Master NODE #################

#bash add_osd_slave.sh IP_ADDRESS_OF_SLAVE

ceph_user="cibadmin"
ceph_password="cibadmin"
user_home="/home/$ceph_user"
CEPH_LOG="/var/log/megam/megamcib/ceph.log"

echo "add_osd_master.sh start execution ====>" >> $CEPH_LOG
echo "Remote osd ip $1 ====>" >> $CEPH_LOG

sudo apt-get -y install sshpass ntp >> $CEPH_LOG

sshpass -p "$ceph_password" scp -o StrictHostKeyChecking=no /home/$ceph_user/.ssh/id_rsa.pub $ceph_user@$1:/home/$ceph_user/.ssh/authorized_keys


REMHOST='ssh $ceph_user@$1 hostname'

#add entry to /etc/hosts
echo "$1 $REMHOST" >> /etc/hosts

remote_osd=`ssh $ceph_user@$1 ls /storage*`
ceph_osds=""
for d in $remote_osd ; do
	ceph_osds="$ceph_osds $REMHOST:$d/osd"
done

ceph-deploy --overwrite-conf osd prepare $ceph_osds >> $CEPH_LOG
ceph-deploy --overwrite-conf osd activate $ceph_osds >> $CEPH_LOG

#IN FIRST HOST
echo "Transfering keys to $1 ====>" >> $CEPH_LOG
scp $user_home/ceph-cluster/ceph.bootstrap-osd.keyring $ceph_user@$1:$user_home/ceph.keyring
scp $user_home/ceph-cluster/*.keyring $ceph_user@$1:$user_home/

ssh $ceph_user@$1 'sudo mv $user_home/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring; sudo mv $user_home/*.keyring /etc/ceph/'


#Libvirt keys
scp $user_home/ceph-cluster/client.libvirt.key $ceph_user@$1:$user_home/ceph-cluster/
scp $user_home/ceph-cluster/secret.xml $ceph_user@$1:$user_home/ceph-cluster/
scp $user_home/ceph-cluster/uid $ceph_user@$1:$user_home/ceph-cluster/

#After this execution ends, run ceph_one_install_slave.sh in slave system.

echo "add_osd_master.sh end execution ====>" >> $CEPH_LOG
