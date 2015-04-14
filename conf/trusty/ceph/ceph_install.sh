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
#
#bash ceph_install.sh install osd1="/storage1" osd2="/storage2" osd3="/storage3"
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
ceph_user="cibadmin"
ceph_password="cibadmin"
ceph_group="cibadmin"
user_home="/home/$ceph_user"

host=`hostname`

osd1="/storage1"
osd2="/storage2"
osd3="/storage3"

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
      osd1=*)
      osd1="${i#*=}"
      ;;
      osd2=*)
      osd2="${i#*=}"
      ;;
      osd3=*)
      osd3="${i#*=}"
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
# Install ceph in trusty or debian
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

install_ceph() {
#ceph user as sudoer 
echo "Make ceph user as sudoer" >> $CEPH_LOG
echo "$ceph_user ALL = (root) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/$ceph_user
#echo "cibadmin ALL = (root) NOPASSWD:ALL" | sudo tee /etc/sudoers.d/cibadmin
sudo chmod 0440 /etc/sudoers.d/$ceph_user
#sudo chmod 0440 /etc/sudoers.d/cibadmin

#Ceph install
echo "Started installing ceph" >> $CEPH_LOG
sudo echo deb http://ceph.com/debian-giant/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list
sudo wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -
sudo apt-get -y update >> $CEPH_LOG
sudo apt-get -y install ceph-deploy ceph-common ceph-mds dnsmasq openssh-server ntp sshpass >> $CEPH_LOG
#sudo apt-get -y install dnsmasq openssh-server ntp sshpass  >> $CEPH_LOG

IP_ADDR=$( getip )

#edit /etc/hosts to access osd nodes
echo "Adding entry in /etc/hosts" >> $CEPH_LOG
echo "$IP_ADDR $host" >> /etc/hosts

echo "Processing ssh-keygen" >> $CEPH_LOG
sudo -u $ceph_user bash << EOF
#Create ssh files
ssh-keygen -N '' -t rsa -f $user_home/.ssh/id_rsa
cp $user_home/.ssh/id_rsa.pub $user_home/.ssh/authorized_keys
EOF

#No prompt on "Add ip to known_hosts list"
sudo -H -u $ceph_user bash -c "cat > /$user_home/.ssh/ssh_config <<EOF
ConnectTimeout 5
Host *
StrictHostKeyChecking no
EOF"

sudo -H -u $ceph_user bash -c "cat > $user_home/.ssh/config <<EOF
Host $host
   Hostname $host
   User $ceph_user
EOF"

echo "Making directory inside osd drive " >> $CEPH_LOG
mkdir $osd1/osd
mkdir $osd2/osd
mkdir $osd3/osd

  #GET first three values of ip
  ip3=`echo $IP_ADDR| cut -d'.' -f 1,2,3`

echo "Ceph configuration started..." >> $CEPH_LOG
sudo -u $ceph_user bash << EOF
mkdir $user_home/ceph-cluster
cd $user_home/ceph-cluster

ceph-deploy new $host

echo "osd crush chooseleaf type = 0" >> ceph.conf
echo "public network = $ip3.0/24" >> ceph.conf
#echo "public network = 7.7.9.0/24" >> ceph.conf
echo "cluster network = $ip3.0/24" >> ceph.conf

#echo "cluster network = 192.168.6.0/24" >> ceph.conf

ceph-deploy install $host
#PROMPT  Are you sure you want to continue connecting (yes/no)? yes    for the first time          cephuser@hostname's password: 

ceph-deploy mon create-initial

#ceph-deploy  --overwrite-conf osd prepare megamubuntu:/storage1/osd megamubuntu:/storage2/osd megamubuntu:/storage3/osd 
#ceph-deploy  --overwrite-conf osd activate megamubuntu:/storage1/osd megamubuntu:/storage2/osd megamubuntu:/storage3/osd
#ceph-deploy osd prepare megamubuntu:/storage1/osd megamubuntu:/storage2/osd megamubuntu:/storage3/osd

ceph-deploy osd prepare $host:$osd1/osd $host:$osd2/osd $host:$osd3/osd
ceph-deploy osd activate $host:$osd1/osd $host:$osd2/osd $host:$osd3/osd

ceph-deploy admin $host
sudo chmod +r /etc/ceph/ceph.client.admin.keyring

sleep 180
ceph osd pool set rbd pg_num 150
#It takes some more time
#better sleep 2 mins
sleep 180
ceph osd pool set rbd pgp_num 150

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
#This function will uninstall ceph
#--------------------------------------------------------------------------

uninstall_ceph() {
ceph-deploy purgedata $host
ceph-deploy forgetkeys
ceph-deploy purge $host
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
