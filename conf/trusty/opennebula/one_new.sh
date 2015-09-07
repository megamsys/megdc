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
# Opennebula Packages install
# This script currently supports ubuntu 14.04 trusty.
#./one.sh install
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
      install_one
      ;;
      [uU][nN][iI][nN][sS][tT][aA][lL][lL])
      uninstall_one
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
  cecho  "Usage    : one.sh [Options]" $green
  cecho  "help     : prints the help message." $blue
  cecho  "install  : Installs opennebula server packages" $blue
  cecho  "uninstall: uninstalls opennebula server packages" $blue
}



ONE_LOG="/var/log/megam/megamcib/one.log"


ping -c 1 get.megam.io &> /dev/null

if [ $? -ne 0 ]; then
	echo "`date`: check your network connection. get.megam.io is down or not reachable!" >> $ONE_LOG
  exit 1
fi

host=`hostname`
echo "Adding entries in /etc/hosts" >> $ONE_LOG

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


#For apt-add-repository command
sudo apt-get -y install software-properties-common python-software-properties >> $ONE_LOG

preinstall(){

}
main_inatll(){

}
postinstall(){

}

install_one(){
preinstall
main_install
postinstall


}


uninstall_one(){


}


