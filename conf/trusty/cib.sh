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
# A linux script which helps to verify the cib installation.
#                      start megam
#                      start one.
###############################################################################

txtblk='\e[0;30m' # Black - Regular
txtred='\e[0;31m' # Red
txtgrn='\e[0;32m' # Green
txtylw='\e[0;33m' # Yellow
txtblu='\e[0;34m' # Blue
txtpur='\e[0;35m' # Purple
txtcyn='\e[0;36m' # Cyan
txtwht='\e[0;37m' # White
bldblk='\e[1;30m' # Black - Bold
bldred='\e[1;31m' # Red
bldgrn='\e[1;32m' # Green
bldylw='\e[1;33m' # Yellow
bldblu='\e[1;34m' # Blue
bldpur='\e[1;35m' # Purple
bldcyn='\e[1;36m' # Cyan
bldwht='\e[1;37m' # White
unkblk='\e[4;30m' # Black - Underline
undred='\e[4;31m' # Red
undgrn='\e[4;32m' # Green
undylw='\e[4;33m' # Yellow
undblu='\e[4;34m' # Blue
undpur='\e[4;35m' # Purple
undcyn='\e[4;36m' # Cyan
undwht='\e[4;37m' # White
bakblk='\e[40m'   # Black - Background
bakred='\e[41m'   # Red
bakgrn='\e[42m'   # Green
bakylw='\e[43m'   # Yellow
bakblu='\e[44m'   # Blue
bakpur='\e[45m'   # Purple
bakcyn='\e[46m'   # Cyan
bakwht='\e[47m'   # White
txtrst='\e[0m'    # Text Reset


CIB_LOG="/var/log/megam/megamcib/cibreport.log"

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
      [rR][eE][pP][oO][rR][tT])
      report_cib
      ;;
      [sS][tT][aA][rR][tT])
      start_cib
      ;;
      [sS][tT][oO][pP])
      stop_cib
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
  echo  -e "${txtgrn}Usage    : ${txtblu}cib.sh [Options]${txtrst}"
  echo  "help     : prints the help message."
  echo  "report   : report about the cib installation"
  echo  "start    : starts cib" $blue
  echo  "stop     : stop cib" $blue
}
#--------------------------------------------------------------------------
# Verify  the cib
#--------------------------------------------------------------------------
report_cib() {
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  echo -e "Select an option to report :"
	echo -e "${bldwht}1)${txtrst} ${txtgrn}megam${txtrst}"
	echo -e "${bldwht}2)${txtrst} ${txtgrn}cobblerd${txtrst}"
	echo -e "${bldwht}3)${txtrst} ${txtgrn}one${txtrst}"
	echo -e "${bldwht}4)${txtrst} ${txtgrn}onehost${txtrst}"
	echo -e "${bldwht}5)${txtrst} ${txtgrn}ceph${txtrst}"
	echo -e "${bldwht}6)${txtrst} ${txtgrn}drbd${txtrst}"
	echo -e "${bldwht}7)${txtrst} ${txtgrn}all${txtrst}"
	read case;

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldylw}%-40s${txtrst}" "   Reporting :";


	case $case in
    1)
		printf "${bldwht}%-10s${txtrst}*\n" "MEGAM";
		report_megam
    ;;
		2)
    printf "${undpur}%-10s${txtrst}\n" "COBBLERD";
		report_cobblerd
		;;
    3)
		printf "${undpur}%-10s${txtrst}\n" "ONE";
		report_one
		;;
		4)
    printf "${undpur}%-10s${txtrst}\n" "ONEHOST";
    report_onehost
		;;
    5)
		printf "${undpur}%-10s${txtrst}\n" "CEPH";
		report_ceph
		;;
    6)
		printf "${undpur}%-10s${txtrst}\n" "DRBD";
		report_drbd
		;;
    7)
    printf "${undpur}%-10s${txtrst}\n" "ALL";
    report_megam
    report_cobblerd
    report_one
    report_onehost
    report_ceph
    report_drbd
    ;;
    *)
    printf "\n"
    ;;
	esac
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";



}


report_megam() {

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Installation : Megam";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  pkgnames=( megamcommon megamcib megamcibn megamnilavu megamsnowflake megamgateway megamd megamchefnative megamanalytics megamdesigner megammonitor riak rabbitmq-server nodejs sqlite3 ruby2.0 openjdk-7-jdk)

  howdy_pkgs pkgnames[@]

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Services : Megam";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  sernames=( megamcib megamcibn megamnilavu snowflake megamgateway megamd megamchefnative megamheka megamanalytics megamdesigner riak )

  howdy_services sernames[@]

  if (( $(ps -ef | grep -v grep | grep "rabbitmq-server" | wc -l) > 0 ))
  then
  printf "${bldpur}%-20s ${bldcyn}%-15s${txtrst} %-5s " 'rabbitmq-server' 'SERVICE'  '.....';
  sudo rabbitmqctl status > /dev/null 2>&1 && {
    printf "${bldgrn}%-15s${txtrst}\n" '[OK]';
  } || {
    printf "${bldred}%-15s${txtrst}\n" '[FAIL]';
  }
  else
  printf "${bldpur}%-15s ${txtred}%-15s${txtrst} %-6s ${bldred}%-15s${txtrst}\n" "rabbitmq-server" 'SERVICE'  '.....'  '[FAIL]';
  fi


}
#--------------------------------------------------------------------------
# Report on cobblerd
#--------------------------------------------------------------------------
report_cobblerd() {
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Installation : Cobblerd";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  pkgnames=( cobbler dnsmasq apache2 debmirror )

  howdy_pkgs pkgnames[@]

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Services : Cobblerd";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  sernames=( xinetd dnsmasq cobbler )

  howdy_services sernames[@]

}
#--------------------------------------------------------------------------
#This function will print out an install report
#--------------------------------------------------------------------------
report_one() {

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Installation : One";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  pkgnames=( opennebula opennebula-sunstone )

  howdy_pkgs pkgnames[@]

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Services : One";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  sernames=( one )

  howdy_services sernames[@]

}

report_ceph() {

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Installation : Ceph";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  pkgnames=( ceph-deploy ceph-common ceph-mds )

  howdy_pkgs pkgnames[@]

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Services : Ceph";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

if [[ `ceph status` == *"HEALTH_OK"* ]]
then
  printf "${bldpur}%-20s ${bldcyn}%-15s${txtrst} %-5s " Ceph 'SERVICE'  '.....'  '[OK]';
  printf "\n";
else
  printf "${bldpur}%-20s ${txtred}%-15s${txtrst} %-5s ${bldred}%-15s${txtrst}\n" ceph 'SERVICE'  '.....'  '[FAIL]';
fi

}


report_drbd() {

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Installation : DRBD";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  pkgnames=( drbd8-utils linux-image-extra-virtual pacemaker heartbeat )

  howdy_pkgs pkgnames[@]

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Services : DRBD";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  sernames=( drbd heartbeat )

  howdy_services sernames[@]

if [[ `crm status` == *"Started"* ]]
then
  printf "${bldpur}%-20s ${bldcyn}%-15s${txtrst} %-5s \n " CRM 'SERVICE'  '.....'  '[OK]';
  printf "\n";
else
  printf "${bldpur}%-20s ${txtred}%-15s${txtrst} %-5s ${bldred}%-15s${txtrst}\n" CRM 'SERVICE'  '.....'  '[FAIL]';
fi

}



#--------------------------------------------------------------------------
# Starts the cib
#--------------------------------------------------------------------------
report_onehost() {
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Installation : One Host";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  pkgnames=( opennebula-node qemu-kvm )

  howdy_pkgs pkgnames[@]

  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  printf "*${bldblu}%-50s${txtrst}*\n" "   Services : One Host";
  printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";

  sernames=( onevm )

  howdy_services sernames[@]

}
#--------------------------------------------------------------------------
# Starts the cib
#--------------------------------------------------------------------------
start_cib() {
  echo -e "${bldylw}Starting cib..${txtrst}"
}
#--------------------------------------------------------------------------
#This function will print out an install report
#--------------------------------------------------------------------------
stop_cib() {
  echo -e "${bldylw}Stopping cib..${txtrst}"
}
#--------------------------------------------------------------------------
#This function will verify if the package exists
#--------------------------------------------------------------------------
howdy_pkgs() {
  pkgnames=("${!1}")
  for pkgname in ${pkgnames[@]}
    do
      dpkg -s "$pkgname" >/dev/null 2>&1 && {
        printf "${bldpur}%-20s ${bldcyn}%-15s${txtrst} %-5s ${bldgrn}%-15s${txtrst}\n" $pkgname 'INSTALL'  '.....' '[OK]';
    } || {
      printf "${bldpur}%-20s ${bldcyn}%-15s${txtrst} %-5s ${bldred}%-15s${txtrst}\n" $pkgname 'INSTALL'    '.....' '[FAIL]';
    }
  done
}
#--------------------------------------------------------------------------
#This function will verify if a process is running, and an upstart service is running
#--------------------------------------------------------------------------
howdy_services(){
  sernames=("${!1}")

  for sername in ${sernames[@]}
  do
    if (( $(ps -ef | grep -v grep | grep $sername | wc -l) > 0 ))
    then
    printf "${bldpur}%-20s ${bldcyn}%-15s${txtrst} %-5s " $sername 'SERVICE'  '.....';
    sudo service $sername status > /dev/null 2>&1 && {
      printf "${bldgrn}%-15s${txtrst}\n" '[OK]';
    } || {
      printf "${bldred}%-15s${txtrst}\n" '[FAIL]';
    }
    else
    printf "${bldpur}%-20s ${txtred}%-15s${txtrst} %-5s ${bldred}%-15s${txtrst}\n" $sername 'SERVICE'  '.....'  '[FAIL]';
    fi
  done

}
#--------------------------------------------------------------------------
#This function will exit out of the script.
#--------------------------------------------------------------------------
exitScript(){
  exit $@
}

#parse parameters
parseParameters "$@"

echo "Good bye."
exitScript 0
