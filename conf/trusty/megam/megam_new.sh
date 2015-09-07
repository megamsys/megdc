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
# Megam Packages install
# This script currently supports ubuntu 14.04 trusty.
#./megam.sh install megamnilavu megamgateway megamd
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

    case $1 in
      [hH][eE][lL][pP])
      help
      ;;
      ('/?')
      help
      ;;
      [iI][nN][sS][tT][aA][lL][lL])
      install_megam
      ;;
      [uU][nN][iI][nN][sS][tT][aA][lL][lL])
      uninstall_megam
      ;;
      *)
      cecho "Unknown option : $item - refer help." $red
      help
      ;;
    esac



  for item in "$@"
  do
    index=$(($index+1))
  done
}
#--------------------------------------------------------------------------
#prints the help to out file.
#--------------------------------------------------------------------------
help() {
  cecho  "Usage    : megam.sh [Options]" $green
  cecho  "help     : prints the help message." $blue
  cecho  "install  : Installs listed megam packages. Pass PACKAGE_NAMES(space seperated)" $blue
  cecho  "uninstall: uninstalls listed megam packages. Pass PACKAGE_NAMES(space seperated)" $blue
}



MEGAM_LOG="/var/log/megam/megamcib/megam.log"


ping -c 1 get.megam.io &> /dev/null

if [ $? -ne 0 ]; then
	echo "`date`: check your network connection. get.megam.io is down or not reachable!" >> $MEGAM_LOG
  exit 1
fi

host=`hostname`
echo "Adding entries in /etc/hosts" >> $MEGAM_LOG

get_ip(){
	while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route
}
get_ip

#ADD /etc/hosts entries
echo "127.0.0.1 `hostname` localhost" >> /etc/hosts
echo "$ipaddr `hostname` localhost" >> /etc/hosts
echo "/etc/hosts entries added"  >> $MEGAM_LOG

#For apt-add-repository command
sudo apt-get -y install software-properties-common python-software-properties >> $MEGAM_LOG

add-apt-repository 'deb [arch=amd64] http://get.megam.io/0.9/ubuntu/14.04/ testing megam'
apt-key adv --keyserver keyserver.ubuntu.com --recv B3E0C1B7
apt-get -y update
apt-get -y install megamcommon


install_megam(){

  for item in "$@"
  do
    index=$(($index+1))

    case $item in
      megamnilavu)
      apt-get install megamnilavu
      ;;
      *)
      cecho "Unknown option : $item - refer help." $red
      help
      ;;
    esac

  done

}


uninstall_megam(){

  for item in "$@"
  do
    index=$(($index+1))

    case $item in
      megamnilavu)
      apt-get remove megamnilavu
      apt-get purge megamnilavu
      ;;
      *)
      cecho "Unknown option : $item - refer help." $red
      help
      ;;
    esac

  done

}


