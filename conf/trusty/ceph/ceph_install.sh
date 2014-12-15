#!/bin/bash
#Copyright (c) 2014 Megam Systems.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
###############################################################################
# A ceph linux script which sets up a ceph-mon and two ceph-osd
# This script currently supports ubuntu 14.04 trusty.
###############################################################################

black='\033[30m'
red='\033[31m'
green='\033[32m'
yellow='\033[33m'
blue='\033[34m'
magenta='\033[35m'
cyan='\033[36m'
white='\033[37m'

alias Reset="tput sgr0"      #  Reset text attributes to normal
# without clearing screen.

CEPH_LOG="/var/log/megam/megamcib/ceph.log"
ceph_user="megamceph"
ceph_password="megamceph"
ceph_group="megamceph"
user_home="/home/$ceph_user"

osd1_ip="$2"
osd1_host="$3"

osd2_ip="$4"
osd2_host="$5"

osd3_ip="$6"
osd3_host="$7"

host=`hostname`

#--------------------------------------------------------------------------
#colored echo
# Argument $1 = message
# Argument $2 = color (
#--------------------------------------------------------------------------
cecho () {
  local default_msg="No message passed."  # Doesn't really need to be a local variable.
  message=${1:-$default_msg}              # Defaults to default message.
  color=${2:-$black}                      # Defaults to black, if not specified.
  echo "$color$message"
  Reset                                   # Reset to normal.
  return
}
#--------------------------------------------------------------------------
#parse the input parameters.
# Pattern in case statement is explained below.
# a*)  The letter a followed by zero or more of any
# *a)  The letter a preceded by zero or more of any
#--------------------------------------------------------------------------
parseParameters()   {
  #integer index=0

  if [ $# -lt 1 ]
    then
    help
    exitScript 1
  fi


  for item in "$@"
  do
    case $item in
      [hH][eE][lL][pP])
      help
      ;;
      ('/?')
      help
      ;;
      [iI][nN][sS][tT][aA][lL][lL])
      install_ceph
      ;;
      [uU][nN][iI][nN][sS][tT][aA][lL][lL])
      uninstall_ceph
      ;;
      *)
      cecho "Unknown option : $item - refer help." $red
      help
      ;;
    esac
    index=$(($index+1))
  done
}
#--------------------------------------------------------------------------
#prints the help to out file.
#--------------------------------------------------------------------------
help() {
  cecho  "Usage    : ceph.sh [Options]" $green
  cecho  "help     : prints the help message." $blue
  cecho  "install osd1_ip osd1_host osd2_ip osd2_host : setup ceph" $blue
  cecho  "uninstall: uninstalls ceph and remove the setup" $blue
}
#--------------------------------------------------------------------------
# Install cobbler in trusty or debian
#--------------------------------------------------------------------------
getip(){
while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route
echo $ipaddr
}

create_cephgroup() {
    if ! getent group $ceph_group > /dev/null 2>&1; then
        addgroup --system $ceph_group
    fi
}

create_cephuser() {
    if ! getent passwd $ceph_user > /dev/null 2>&1; then
        useradd -d $user_home -m -g $ceph_group $ceph_user -s /bin/bash
        #Set password
        sudo echo -e "$ceph_password\n$ceph_password\n" | sudo passwd $ceph_user
    else
        user_home=`getent passwd $ceph_user | cut -f6 -d:`
        # Renable user (give him a shell)
        usermod --shell /bin/bash $ceph_user
        # Make sure MEGAMHOME exists, might have been removed on previous purge
        mkdir -p $user_home
    fi
}

install_ceph() {
 
create_cephgroup
create_cephuser
#ceph user as sudoer 
echo "$ceph_user ALL = (root) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/$ceph_user
sudo chmod 0440 /etc/sudoers.d/$ceph_user

#Ceph install
echo deb http://ceph.com/debian-giant/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
sudo apt-get -y update
sudo apt-get -y install ceph-deploy ceph-common ceph-mds

IP_ADDR=$( getip )

sudo apt-get install ntp

#edit /etc/hosts to access osd nodes

echo "$IP_ADDR $host" >> /etc/hosts
echo "$osd1_ip $osd1_host" >> /etc/hosts
echo "$osd2_ip $osd2_host" >> /etc/hosts
echo "$osd3_ip $osd3_host" >> /etc/hosts


sudo apt-get install dnsmasq 
sudo apt-get install openssh-server

sudo -u $ceph_user bash << EOF
#Create ssh files
ssh-keygen -N '' -t rsa -f $user_home/.ssh/id_rsa

sshpass -p "$ceph_password" scp -o StrictHostKeyChecking=no $user_home/.ssh/id_rsa.pub $ceph_user@$osd1_host:$user_home/.ssh/authorized_keys
sshpass -p "$ceph_password" scp -o StrictHostKeyChecking=no $user_home/.ssh/id_rsa.pub $ceph_user@$osd2_host:$user_home/.ssh/authorized_keys
sshpass -p "$ceph_password" scp -o StrictHostKeyChecking=no $user_home/.ssh/id_rsa.pub $ceph_user@$osd3_host:$user_home/.ssh/authorized_keys

sshpass -p "oneadmin" scp -o StrictHostKeyChecking=no /var/lib/one/.ssh/id_rsa.pub oneadmin@192.168.2.14:/var/lib/one/.ssh/authorized_keys
EOF

cat > $user_home/.ssh/config <<EOF
Host $osd1_host
   Hostname $osd1_host
   User $ceph_user
Host $osd2_host
   Hostname $osd2_host
   User $ceph_user
Host $osd3_host
   Hostname $osd3_host
   User $ceph_user
Host $host
   Hostname $host
   User $ceph_user
EOF

chown $ceph_user:$ceph_user $user_home/.ssh/config

mkdir /storage

sudo -u $ceph_user bash << EOF
mkdir $user_home/ceph-cluster
cd $user_home/ceph-cluster

ceph-deploy new $host


ceph-deploy install $osd1_host $osd2_host $osd3_host $host
#PROMPT  Are you sure you want to continue connecting (yes/no)? yes    for the first time          cephuser@hostname's password: 

ceph-deploy mon create-initial

scp $user_home/ceph-cluster/ceph.bootstrap-osd.keyring $ceph_user@$osd1_host:$user_home/ceph.keyring
scp $user_home/ceph-cluster/ceph.bootstrap-osd.keyring $ceph_user@$osd2_host:$user_home/ceph.keyring
scp $user_home/ceph-cluster/ceph.bootstrap-osd.keyring $ceph_user@$osd3_host:$user_home/ceph.keyring

#scp /home/megamceph/ceph-cluster/ceph.bootstrap-osd.keyring megamceph@alrin:/home/megamceph/ceph.keyring

#ssh megamceph@osd1 'sudo mv /home/megamceph/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'


ssh $ceph_user@$osd1_host 'sudo mv $user_home/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'
ssh $ceph_user@$osd2_host 'sudo mv $user_home/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'

ssh $ceph_user@$osd3_host 'sudo mv $user_home/ceph.keyring /var/lib/ceph/bootstrap-osd/; sudo chmod 600 /var/lib/ceph/bootstrap-osd/ceph.keyring'

#ceph-deploy osd prepare megamubuntu:/storage alrin:/storage osd1:/storage
ceph-deploy osd prepare $osd1_host:/storage $osd2_host:/storage $osd3_host:/storage
ceph-deploy osd activate $osd1_host:/storage $osd2_host:/storage $osd3_host:/storage

ceph-deploy admin $host $osd1_host $osd2_host $osd3_host
sudo chmod +r /etc/ceph/ceph.client.admin.keyring

EOF

install_complete

}



install_complete() {
  cecho "##################################################" $green
  cecho "Ceph installed successfully." $yellow
  cecho "##################################################" $green
  echo  "Ceph installed successfully." >> $CEPH_LOG
}

#--------------------------------------------------------------------------
#This function will uninstall cobblerd
#--------------------------------------------------------------------------

uninstall_ceph() {
ceph-deploy purgedata $host $osd1_host $osd2_host $osd3_host
ceph-deploy forgetkeys
ceph-deploy purge $host $osd1_host $osd2_host $osd3_host
sudo rm -r /var/lib/ceph/
sudo apt-get -y autoremove
sudo apt-get -y remove ceph-deploy ceph-common ceph-mds
sudo apt-get -y purge ceph-deploy ceph-common ceph-mds
sudo apt-get -y autoremove

sudo rm -r /run/ceph
sudo rm -r /var/lib/ceph
sudo rm /var/log/upstart/ceph*
sudo rm ~/ceph-cluster/*
}
#--------------------------------------------------------------------------
#This function will exit out of the script.
#--------------------------------------------------------------------------
exitScript(){
  exit $@
}

#parse parameters
parseParameters "$@"

cecho "Good bye." $yellow
exitScript 0
