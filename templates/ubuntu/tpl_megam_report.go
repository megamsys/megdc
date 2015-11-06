/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package ubuntu

import (
	"fmt"
	"github.com/dynport/urknall"
	"github.com/megamsys/megdc/templates"
	"os/exec"
	// "sync"
	"strings"
	"os"
	"path"
	//"reflect"
	//pp "github.com/megamsys/libgo/cmd"
	//"github.com/codeskyblue/go-sh"
)

var ubuntumegamreport *UbuntuMegamReport

func init() {
	ubuntumegamreport = &UbuntuMegamReport{}
	templates.Register("UbuntuMegamReport", ubuntumegamreport)
}

type UbuntuMegamReport struct{}

func (tpl *UbuntuMegamReport) Render(p urknall.Package) {
	p.AddTemplate("report", &UbuntuMegamReportTemplate{})
}

func (tpl *UbuntuMegamReport) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamReport{})
}

type UbuntuMegamReportTemplate struct{}

func (m *UbuntuMegamReportTemplate) Render(pkg urknall.Package) {
	
	if err := writefile(); err != nil {
		return
	} else {
		commands := "bash /var/lib/megam/report.sh megam"

	commandWords := strings.Fields(commands)
	out, err := exe_cmd(commandWords[0], commandWords[1:])
	if err != nil {
		return
	}
	fmt.Println(out)
	}
}

func exe_cmd(cmd string, args []string) (string, error) {
	cmd_out := exec.Command(cmd, args...)
	out, err := cmd_out.Output()

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return string(out), nil
}

func writefile() error {
	basePath := "/var/lib"
	dir := path.Join(basePath, "megam")

	filePath := path.Join(dir, "report.sh")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if errm := os.MkdirAll(dir, 0777); errm != nil {
			return errm
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(string(reportscript)); err != nil {
		return err
	}
	return nil
}

const reportscript = `#!/bin/bash
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
      [mM][eE][gG][aA][mM])
      report_megam
      ;;
      [oO][nN][eE])
      report_one
      ;;
      [oO][nN][eE][hH][oO][sS][tT])
      report_one_host
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
  echo  "megam    : report about the megam packages installation"
  echo  "one      : report about the one installation"
  echo  "one_host     : report about the one_host installation"
}

report_megam() { 

  pkgnames=( megamcommon megamnilavu megamsnowflake megamgateway megamd riak rabbitmq-server ruby2.0 openjdk-7-jdk)

  howdy_pkgs pkgnames[@]  

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
        printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  	printf "*${bldblu}%-50s${txtrst}*\n" "   ${pkgname} Report";
  	printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
        dpkg -s "$pkgname" >/dev/null 2>&1 && {        
           contentparse $pkgname "Status" "Install-"
	   contentparse $pkgname "Size"
	   contentparse $pkgname "Architecture"
	   contentparse $pkgname "Version"
	   contentparse $pkgname "Depends"
	   contentparse $pkgname "License"
	   contentparse $pkgname "Homepage"
	   statusparse  $pkgname "Running-status"
       } || {
        printf "${bldylw}%-20s ${bldylw}%-15s${bldylw}${txtrst} \n" $pkgname 'Not Install';
      }
     printf "*${txtblu}%-50s${txtrst}*\n" "--------------------------------------------------";
  done
}
#--------------------------------------------------------------------------
#This function will verify if a process is running, and an upstart service is running
#--------------------------------------------------------------------------
statusparse(){
  sername=$1
  
   if (( $(ps -ef | grep -v grep | grep $sername | wc -l) > 0 ))
    then
    sudo service $sername status > /dev/null 2>&1 && {
      echo -e "${bldpur}$2${txtrst} \t \t ${bakgrn}Running${txtrst}\n";     
    } || {
      echo -e "${bldpur}$2${txtrst} \t \t ${bakred}NotRunning${txtrst}\n";  
    }
    else
      echo -e "${bldpur}$2${txtrst} \t \t ${bakred}NotRunning${txtrst}\n"; 
    fi
}

contentparse(){
    str=$(dpkg -s $1 | grep $2)
    IFS=: read typename value <<< $str		
    echo -e "${bldpur}$3$typename \t \t ${bldcyn}$value${txtrst}";
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
`
