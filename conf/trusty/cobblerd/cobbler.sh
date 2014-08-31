#!/bin/sh
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
# A cobbler linux script which sets up cobblerd and a DHCP using dnsmasq.
# The dhcp range used by cobbler is 192.168.2.20 - 200.
# The i/p address of the cobbler is 192.168.2.3
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

COBBLER_LOG="/var/log/megam/megamcib/cobbler.log"

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
      install_cobbler
      ;;
      [uU][nN][iI][nN][sS][tT][aA][lL][lL])
      uninstall_cobbler
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
  cecho  "Usage    : cobbler.sh [Options]" $green
  cecho  "help     : prints the help message." $blue
  cecho  "install  : setup cobblerd which network boots on dhcp 192.168.2.x" $blue
  cecho  "uninstall: uninstalls cobblerd and remove the setup" $blue
}
#--------------------------------------------------------------------------
# Install cobbler in trusty or debian
#--------------------------------------------------------------------------
install_cobbler() {
  ping -c 1 us.archive.ubuntu.com &> /dev/null

  if [ $? -ne 0 ]; then
    echo "`date`: check your network connection. us.archive.ubuntu.com is not reachable!" >> $COBBLER_LOG
    exit 1
  fi

  cecho "Installing cobblerd.." $yellow
  apt-get -y install cobbler cobbler-common dhcp3-server  >> $COBBLER_LOG
  apt-get -y debmirror >> $COBBLER_LOG
  cobbler get-loaders >> $COBBLER_LOG
  cobbler check >> $COBBLER_LOG
   
  #If any errors just fix as it says.
  #on success --> No configuration problems found.  All systems go.
  cobbler sync

  #this can't be done't here as we are in a background mode now.
  dpkg-reconfigure cobbler $COBBLER_LOG

  configure_cobbler
}
#--------------------------------------------------------------------------
#Configures cobbler with the dhcp range 192.168.2.20-200.
#cobblerd machine ip address is 192.168.2.3
#dhcp managed by dnsmasq
#--------------------------------------------------------------------------
configure_cobbler() {
  cecho "Configuring cobblerd.." $yellow
  
  echo 'base_megamreporting_enabled: 1' >> /etc/cobbler/settings

  sed -i 's/manage_dhcp: 0/manage_dhcp: 1/g' /etc/cobbler/settings

  cp install_post_cibnode.py /usr/lib/python2.7/cobbler/modules

  service cobbler restart

  cobbler sync

  echo "manage_dhcp:1 => dhcp managed by cobbler.."

  apt-get install xinetd tftpd tftp >> $COBBLER_LOG

  sed -i 's/^[ \t]*option routers.*/option routers 192.168.2.3;/' /etc/cobbler/dhcp.template
  echo "route:192.168.2.3 => route is your cobblerd machine.."

  sed -i 's/^[ \t]*option domain-name-servers.*/option domain-name-servers 192.168.2.3;/' /etc/cobbler/dhcp.template
  echo "domain-name-servers 192.168.2.3 => domain-name-servers is your cobbled machine.."

  sed -i 's/^[ \t]*option subnet-mask.*/option subnet-mask 255.255.255.0;/' /etc/cobbler/dhcp.template
  echo "manage_dhcp:1 dhcp managed by cobbler.."

  sed -i 's/^[ \t]*range dynamic-bootp.*/range dynamic-bootp 192.168.2.20 192.168.2.100;/' /etc/cobbler/dhcp.template
  echo "manage_dhcp:1 dhcp managed by cobbler.."

  sed -i 's/module = manage_bind/module = manage_dnsmasq/g' /etc/cobbler/modules.conf
  echo "manage_bind:manage_dnsmasq => use dnsmasq.."

  sed -i 's/module = manage_isc/module = manage_dnsmasq/g' /etc/cobbler/modules.conf
  echo "manage_isc:manage_dnsmanq => use dnsmasq.."

  sed -i 's/^[ \t]*dhcp-range=.*/dhcp-range=192.168.2.20,192.168.2.200/' /etc/dnsmasq.conf
  echo "dhcp-range=192.168.2.20-200 => use dhcp range from 192.168.2-200.." $blue

  sed -i 's/^[ \t]*dhcp-option=3.*/dhcp-option=3,192.168.2.23/' /etc/dnsmasq.conf

  echo "enable-tftp" >> /etc/dnsmasq.conf
  echo "enable-tftp.."

  echo "tftp-root=/var/lib/tftpboot" >> /etc/dnsmasq.conf

  sed -i 's/^[ \t]*user.*/user                    = root/' /etc/cobbler/tftpd.template
  sed -i 's/^[ \t]*server_args.*/server_args             = -v -s /var/lib/tftpboot/' /etc/cobbler/tftpd.template

  service xinetd restart

  service dnsmasq restart

  service cobbler restart

  cobbler sync

  setup_boottrusty
}
#--------------------------------------------------------------------------
#Download and setup trusty amd64 mini
#--------------------------------------------------------------------------
setup_boottrusty() {
  cecho "Setup trusty megam node.." $yellow

  cd /var/lib/cobbler/isos/

  wget --tries=3 -c https://s3-ap-southeast-1.amazonaws.com/megampub/iso/trusty_megamnode.iso

  mv trusty_megamnode.iso trusty-amd64-megamnode.iso

  mount -o loop trusty-amd64-megamnode.iso /mnt

  cobbler import --name=ubuntu-server-trusty-megamnode --path=/mnt --breed=ubuntu --arch x86_64

  # Do you have to unmount the /mnt directory ?  after you are done importing ?

  service xinetd restart

  service dnsmasq restart

  service cobbler restart

  cobbler sync

  install_complete
}
#--------------------------------------------------------------------------
#This function will print out an install report
#--------------------------------------------------------------------------
install_complete() {
  cecho "##################################################" $green
  cecho "Step 1: cobbler installed successfully." $yellow
  cecho "        The ip address of cobblerd is 192.168.2.3"
  cecho "        The  subnet   dhcp range   is [192.168.2.20 .. 200]"
  cecho "Refer http://bit.ly/megamcib for more information." $yellow
  cecho "##################################################" $green
}
#--------------------------------------------------------------------------
#This function will uninstall cobblerd
#--------------------------------------------------------------------------
uninstall_cobbler() {
  cecho "Uninstalling cobblerd.." $yellow


  apt-get -y remove cobbler cobbler-common cobbler-web dhcp3-server >> $COBBLER_LOG

  apt-get -y remove debmirror >> $COBBLER_LOG

  apt-get -y purge cobbler cobbler-common cobbler-web dhcp3-server >> $COBBLER_LOG

  apt-get -y purge debmirror >> $COBBLER_LOG

  apt-get -y remove xinetd tftpd tftp >> $COBBLER_LOG

  apt-get -y purge xinetd tftpd tftp >> $COBBLER_LOG

  cecho "##################################################" $green

  cecho "Uninstall complete.." $yellow
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
