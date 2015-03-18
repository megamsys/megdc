#!/bin/bash
########### IN Second NODE #################
echo deb http://ceph.com/debian-giant/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
sudo apt-get -y update
sudo apt-get -y install ceph-deploy ceph-common ceph-mds


sudo mkdir /storage4/osd
sudo mkdir /storage5/osd

#ceph osd pool set rbd pg_num 150
#It takes some more time
#better sleep 2 mins
#ceph osd pool set rbd pgp_num 150

#echo "public network = 192.168.1.0/24" >> ceph.conf
#echo "cluster network = 192.168.1.0/24" >> ceph.conf

#######MON#######
#ceph-deploy --overwrite-conf mon create megamslave

########### IN FIRST NODE #################

#osd_number=`ceph osd create`
sudo apt-get -y install sshpass
#add entry to /etc/hosts
echo "192.168.1.101 megamslave" > /etc/hosts

ceph-deploy --overwrite-conf mon create megamslave
#ceph-deploy --overwrite-conf mon create megamaathi
#ceph-deploy --overwrite-conf mon create megamsaaral


#CREATE directory ~/.ssh in slave system
sshpass -p "cibadmin" scp -o StrictHostKeyChecking=no /home/cibadmin/.ssh/id_rsa.pub cibadmin@megamslave:/home/cibadmin/.ssh/authorized_keys
sshpass -p "cibadmin" scp -o StrictHostKeyChecking=no /home/cibadmin/.ssh/id_rsa.pub cibadmin@megamaathi:/home/cibadmin/.ssh/authorized_keys
sshpass -p "cibadmin" scp -o StrictHostKeyChecking=no /home/cibadmin/.ssh/id_rsa.pub cibadmin@megamsaaral:/home/cibadmin/.ssh/authorized_keys

ceph-deploy --overwrite-conf osd prepare megamslave:/storage4/osd megamslave:/storage5/osd
#ceph-deploy --overwrite-conf osd prepare megamaathi:/storage7/osd megamaathi:/storage8/osd
#ceph-deploy --overwrite-conf osd prepare megamsaaral:/storage9/osd megamsaaral:/storage10/osd
ceph-deploy osd activate megamslave:/storage4/osd megamslave:/storage5/osd




########### IN Second NODE #################
#sudo ln -s /storage4/osd  /var/lib/ceph/osd/ceph-$osd_number

#IN FIRST HOST
scp /home/cibadmin/ceph-cluster/ceph.bootstrap-osd.keyring cibadmin@megamslave:/home/cibadmin/ceph.keyring
#scp /home/cibadmin/ceph-cluster/ceph.bootstrap-osd.keyring cibadmin@megamaathi:/home/cibadmin/ceph.keyring
#scp /home/cibadmin/ceph-cluster/ceph.bootstrap-osd.keyring cibadmin@megamsaaral:/home/cibadmin/ceph.keyring
ssh cibadmin@megamslave 'sudo mv /home/cibadmin/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'
#ssh cibadmin@megamaathi 'sudo mv /home/cibadmin/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'
#ssh cibadmin@megamsaaral 'sudo mv /home/cibadmin/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'
#sudo ceph-osd -i 3 --mkfs --mkkey
#scp /home/cibadmin/ceph-cluster/ceph.conf cibadmin@megamslave:/home/cibadmin/ceph.conf
#ssh cibadmin@megamslave 'sudo mv /home/cibadmin/ceph.conf /etc/ceph/; sudo chmod 600 /etc/ceph/ceph.conf'

#ceph osd crush set osd.3 3 3.0 pool=default rack=unknownrack host=megamslave

#ceph auth add osd.{osd-num} osd 'allow *' mon 'allow rwx' -i /var/lib/ceph/osd/ceph-{osd-num}/keyring
#ceph auth add osd.3 osd 'allow *' mon 'allow rwx' -i /var/lib/ceph/osd/ceph-3/keyring

#ceph-osd -i $osd_number --mkfs --mkkey



#scp /home/cibadmin/ceph-cluster/ceph.bootstrap-osd.keyring cibadmin@megamslave:/home/cibadmin/ceph.keyring

scp /etc/ceph/*libvirt.keyring cibadmin@megamslave:/home/cibadmin/
scp ./*.keyring cibadmin@megamaathi:/home/cibadmin/
scp ./*.keyring cibadmin@megamsaaral:/home/cibadmin/

In each servers, sudo mv ~/*.keyring /etc/ceph/
