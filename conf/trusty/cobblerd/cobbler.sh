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
  apt-get -y install cobbler cobbler-common cobbler-web dhcp3-server dnsmasq  >> $COBBLER_LOG

  apt-get -y install debmirror >> $COBBLER_LOG
  cobbler get-loaders >> $COBBLER_LOG
  cobbler check >> $COBBLER_LOG

  #If any errors just fix as it says.
  #on success --> No configuration problems found.  All systems go.
  cobbler sync

  install_reconfigure_cobbler

  configure_cobbler
}

#--------------------------------------------------------------------------
#This runs the exact steps the dpkg-reconfigure does for cobbler.
#--------------------------------------------------------------------------
install_reconfigure_cobbler() {
  cecho "dpkg reconfiguring cobblerd.." $yellow
  password="cobbler"
	hash=$(printf "cobbler:Cobbler:$password" | md5sum | awk '{print $1}')
	[ -e /etc/cobbler/users.digest ] || install -o root -g root -m 0600 /dev/null /etc/cobbler/users.digest
	htpasswd -D /etc/cobbler/users.digest "cobbler" || true
	printf "cobbler:Cobbler:$hash\n" >> /etc/cobbler/users.digest
	hash=$(printf "$password" | openssl passwd -1 -stdin)
	sed -i "s%^default_password_crypted:.*%default_password_crypted: \"$hash\"%" /etc/cobbler/settings

	while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route
  cecho "figured out the host ip" $yellow
	echo $ipaddr

	if grep -qs "^next_server: *..*..*..*$" /etc/cobbler/settings; then
		sed -i "s/^next_server: *..*..*..*$/next_server: $ipaddr/" /etc/cobbler/settings
	fi
	if grep -qs "^server: *..*..*..*$" /etc/cobbler/settings; then
		sed -i "s/^server: *..*..*..*$/server: $ipaddr/" /etc/cobbler/settings
	fi

	# Enable required apache modules
	a2enmod proxy_http
	a2enmod wsgi
	a2enmod rewrite

	# Install cobbler files and web config for API
	ln -sf /var/lib/cobbler/webroot/cobbler /var/www/cobbler
	if [ ! -e /etc/apache2/conf-enabled/cobbler_web.conf ]; then
	    ln -sf /etc/cobbler/cobbler.conf /etc/apache2/conf-enabled/cobbler.conf
	fi

	# Need to restart apache to pickup web configs
	if [ -x /usr/sbin/invoke-rc.d ]; then
		invoke-rc.d apache2 restart || true
	else
		/etc/init.d/apache2 restart || true
	fi
}
#--------------------------------------------------------------------------
#Configures cobbler with the dhcp range 192.168.2.20-200.
#cobblerd machine ip address is 192.168.2.3
#dhcp managed by dnsmasq
#--------------------------------------------------------------------------
configure_cobbler() {
  cecho "configuring cobblerd.." $yellow

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

  sed -i 's/^[ \t]*dhcp-range=.*/dhcp-range=192.168.2.20,192.168.2.200/' /etc/cobbler/dnsmasq.template
  sed -i 's/^[ \t]*dhcp-option=3.*/dhcp-option=3,192.168.2.23/' /etc/cobbler/dnsmasq.template

  echo "enable-tftp" >> /etc/dnsmasq.conf
  echo "tftp-root=/var/lib/tftpboot" >> /etc/dnsmasq.conf
  echo "enable-tftp" >> /etc/cobbler/dnsmasq.template
  echo "tftp-root=/var/lib/tftpboot" >> /etc/cobbler/dnsmasq.template
  echo "enable-tftp.."

  echo "tftp-root=/var/lib/tftpboot" >> /etc/dnsmasq.conf

  sed -i 's/^[ \t]*user.*/user                    = root/' /etc/cobbler/tftpd.template
  sed -i 's/^[ \t]*server_args.*/server_args             = -v -s /var/lib/tftpboot/' /etc/cobbler/tftpd.template
  sed -i 's/^[ \t]*server .*/server             = /usr/sbin/in.tftpd/' /etc/cobbler/tftpd.template


  restart_services

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



  restart_services

  cecho "Running Synchronisation..." $yellow

  sleep 3
  cobbler sync
  sleep 10
  cobbler reposync
  boot_menu
  install_complete
}

restart_services() {
  service xinetd restart

  service dnsmasq restart

  service cobbler restart
}
#--------------------------------------------------------------------------
#This function will print out boot menu
#--------------------------------------------------------------------------
boot_menu() {
  cecho "##################################################" $green
  > /var/lib/tftpboot/pxelinux.cfg/default

  cat > //var/lib/tftpboot/pxelinux.cfg/default <<EOF
DEFAULT menu
PROMPT 0
MENU TITLE Megam CIB | http://gomegam.com

LABEL ubuntu-server-trusty-megamnode-x86_64
        kernel /images/ubuntu-server-trusty-megamnode-x86_64/linux
        MENU LABEL ubuntu-server-trusty-megamnode-x86_64
        append initrd=/images/ubuntu-server-trusty-megamnode-x86_64/initrd.gz ksdevice=bootif lang=  locale=en_US priority=critical text  auto-install/enable=true priority=critical url=http://144.76.190.227/npreseed.cfg hostname=ubuntu-server-trusty-megamnode-x8664 domain=local.lan suite=trusty
        ipappend 2
MENU end
EOF
restart_services
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


  apt-get -y remove cobbler cobbler-common cobbler-web dhcp3-server dnsmasq >> $COBBLER_LOG

  apt-get -y purge cobbler cobbler-common cobbler-web dhcp3-server >> $COBBLER_LOG

  apt-get -y remove debmirror >> $COBBLER_LOG

  apt-get -y purge debmirror >> $COBBLER_LOG

  apt-get -y remove xinetd tftpd tftp >> $COBBLER_LOG

  apt-get -y purge xinetd tftpd tftp >> $COBBLER_LOG

  [ -d /var/log/cobbler ] && rm -rf /var/log/cobbler
  [ -d /var/lib/cobbler ] && rm -rf /var/lib/cobbler
  [ -d /etc/cobbler ] && rm -rf /etc/cobbler

  apt-get -y autoremove

  rm /var/cache/apt/archives/cobbler-common_2.4.1-0ubuntu2_all.deb
  rm /var/cache/apt/archives/cobbler-web_2.4.1-0ubuntu2_all.deb

  rm /var/cache/apt/archives/cobbler_2.4.1-0ubuntu2_all.deb
  rm /var/log/upstart/cobbler.log

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
