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
# The dhcp range used by cobbler is x.x.x.20 - 200.
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
  cecho  "install  : setup cobblerd which network boots on dhcp 192.168.x.x" $blue
  cecho  "uninstall: uninstalls cobblerd and remove the setup" $blue
}
#--------------------------------------------------------------------------
# Install cobbler in trusty or debian
#--------------------------------------------------------------------------

#--------------------------------------------------------------------------
#This function will restart the services
#--------------------------------------------------------------------------
restart_all() {
  echo "Restarting services ..." >> $COBBLER_LOG
  service xinetd restart
  service dnsmasq restart
  service cobbler restart
  service apache2 restart
  echo "Restarted services..." >> $COBBLER_LOG
}

install_cobbler() {
  ping -c 1 us.archive.ubuntu.com &> /dev/null

  if [ $? -ne 0 ]; then
    echo "`date`: check your network connection. us.archive.ubuntu.com is not reachable!" >> $COBBLER_LOG
    exit 1
  fi

  echo "Installing cobbler packages" >> $COBBLER_LOG
  cecho "Installing cobbler packages.." $yellow

  apt-get -y install cobbler cobbler-common cobbler-web dnsmasq debmirror xinetd tftpd tftp >> $COBBLER_LOG

  cobbler get-loaders >> $COBBLER_LOG
  cobbler check >> $COBBLER_LOG

  echo "Manually dpkgconfigure cobbler" >> $COBBLER_LOG

  manual_dpkgreconfigure

  echo "Configuring cobbler" >> $COBBLER_LOG

  configure_cobbler

  import_trustynode_iso

  setup_profile_trustynode
  #setup_profile_cephnode
  cecho "Running Synchronisation..." $yellow
  echo "Running Synchronisation..." >> $COBBLER_LOG
  sleep 3
  cobbler sync
  sleep 10
  mv /tftpboot/ /var/lib/                       #Who place the file here?
  restart_all
  echo "Syncing repositories..." >> $COBBLER_LOG
  cobbler reposync &
  echo "Reposync completed..." >> $COBBLER_LOG
  restart_all

  install_complete

}
#--------------------------------------------------------------------------
#This runs the exact steps the dpkg-reconfigure does for cobbler.
#--------------------------------------------------------------------------
manual_dpkgreconfigure() {
  cecho "dpkg reconfiguring cobblerd.." $yellow

  password="cobbler"
  hash=$(printf "cobbler:Cobbler:$password" | md5sum | awk '{print $1}')
  [ -e /etc/cobbler/users.digest ] || install -o root -g root -m 0600 /dev/null /etc/cobbler/users.digest
  htpasswd -D /etc/cobbler/users.digest "cobbler" || true
  printf "cobbler:Cobbler:$hash\n" >> /etc/cobbler/users.digest
  hash=$(printf "$password" | openssl passwd -1 -stdin)

  sed -i "s%^default_password_crypted:.*%default_password_crypted: \"$hash\"%" /etc/cobbler/settings

  cecho "reconfigured password"

  configure_with_ip

  # Enable required apache modules
  a2enmod proxy_http
  a2enmod wsgi
  a2enmod rewrite

  # Install cobbler files and web config for API
  ln -sf /var/lib/cobbler/webroot/cobbler /var/www/cobbler
  if [ ! -e /etc/apache2/conf-enabled/cobbler_web.conf ]; then
     ln -sf /etc/cobbler/cobbler.conf /etc/apache2/conf-enabled/cobbler.conf
  fi

  echo "Reconfigure complete." >> $COBBLER_LOG

}
#--------------------------------------------------------------------------
# Figure out the ip address and set it up in the
# /etc/cobbler/settings : server, next_server
# update the ipaddress in /etc/hosts file.
# not yet done : change network to use static, if it uses dhcp (https://github.com/megamsys/cloudinabox/issues/51).
#-------------------------------------------------------------------------
configure_with_ip() {
	while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route

echo "127.0.0.1 localhost" >> /etc/hosts
echo "$ipaddr megamubuntu" >> /etc/hosts

  	  
	  
    if grep -qs "^next_server: *..*..*..*$" /etc/cobbler/settings; then
        sed -i "s/^next_server: *..*..*..*$/next_server: $ipaddr/" /etc/cobbler/settings
    fi

    if grep -qs "^server: *..*..*..*$" /etc/cobbler/settings; then
        sed -i "s/^server: *..*..*..*$/server: $ipaddr/" /etc/cobbler/settings
    fi

	echo "Configured the ip" >> $COBBLER_LOG
}
#--------------------------------------------------------------------------
#Configures cobbler with the dhcp range 192.168.x.20-200.
#cobblerd machine ip address is $ipaddr
#dhcp managed by dnsmasq
#--------------------------------------------------------------------------
configure_cobbler() {
  cecho "configuring cobblerd.." $yellow
  echo "configuring cobblerd.." >> $COBBLER_LOG
  echo 'base_megamreporting_enabled: 1' >> /etc/cobbler/settings

  sed -i 's/manage_dhcp: 0/manage_dhcp: 1/g' /etc/cobbler/settings

  sed -i "s/^[ \t]*option routers.*/option routers $ipaddr;/" /etc/cobbler/dhcp.template
  echo "route:$ipaddr => route is your cobblerd machine.." >> $COBBLER_LOG

  sed -i "s/^[ \t]*option domain-name-servers.*/option domain-name-servers $ipaddr;/" /etc/cobbler/dhcp.template
  echo "domain-name-servers $ipaddr => domain-name-servers is your cobbled machine.." >> $COBBLER_LOG

  sed -i 's/^[ \t]*option subnet-mask.*/option subnet-mask 255.255.255.0;/' /etc/cobbler/dhcp.template
  echo "manage_dhcp:1 dhcp managed by cobbler.." >> $COBBLER_LOG

  #GET first three values of ip
  ip3=`echo $ipaddr| cut -d'.' -f 1,2,3`

  sed -i "s/^[ \t]*range dynamic-bootp.*/range dynamic-bootp $ip3.20 $ip3.100;/" /etc/cobbler/dhcp.template
  echo "manage_dhcp:1 dhcp managed by cobbler.." >> $COBBLER_LOG

  sed -i 's/module = manage_bind/module = manage_dnsmasq/g' /etc/cobbler/modules.conf
  echo "manage_bind:manage_dnsmasq => use dnsmasq.." >> $COBBLER_LOG

  sed -i 's/module = manage_isc/module = manage_dnsmasq/g' /etc/cobbler/modules.conf
  echo "manage_isc:manage_dnsmanq => use dnsmasq.." >> $COBBLER_LOG

  sed -i "s/^[ \t]*dhcp-range=.*/dhcp-range=$ip3.20,$ip3.200/" /etc/dnsmasq.conf
  echo "dhcp-range=$ip3.20-200 => use dhcp range from $ip3.200.."  >> $COBBLER_LOG

  sed -i "s/^[ \t]*dhcp-option=3.*/dhcp-option=3,$ip3.1/" /etc/dnsmasq.conf

  sed -i "s/^[ \t]*dhcp-range=.*/dhcp-range=$ip3.20,$ip3.200/" /etc/cobbler/dnsmasq.template
  sed -i "s/^[ \t]*dhcp-option=3.*/dhcp-option=3,$ip3.1/" /etc/cobbler/dnsmasq.template

  echo "enable-tftp" >> /etc/dnsmasq.conf
  echo "tftp-root=/var/lib/tftpboot" >> /etc/dnsmasq.conf
  echo "enable-tftp" >> /etc/cobbler/dnsmasq.template
  echo "tftp-root=/var/lib/tftpboot" >> /etc/cobbler/dnsmasq.template

  echo "Setup complete to manage DHCP and DNS by cobbler" >> $COBBLER_LOG

  echo "Setup complete to for tftp" >> $COBBLER_LOG

  sed -i 's/^[ \t]*user.*/user                    = root/' /etc/cobbler/tftpd.template
  sed -i "s/^[ \t]*server  .*/server             = \/usr\/sbin\/n.tftpd\//" /etc/cobbler/tftpd.template
  sed -i "s/^[ \t]*server_args.*/server_args             = -v -s \/var\/lib\/tftpboot\//" /etc/cobbler/tftpd.template

}
#--------------------------------------------------------------------------
#Download and setup trusty amd64 mini
#--------------------------------------------------------------------------
import_trustynode_iso() {
  cecho "Setup trusty megam node.." $yellow
  echo "Setting up boot: trusty megam node iso" >> $COBBLER_LOG

  cd /var/lib/cobbler/isos/

  echo "Downloading trusty_megamnode.iso" >> $COBBLER_LOG

  wget --tries=5 -c https://s3-ap-southeast-1.amazonaws.com/megampub/iso/trusty_megamnode.iso

  mv trusty_megamnode.iso trusty-amd64-megamnode.iso

  mount -o loop trusty-amd64-megamnode.iso /mnt

  cobbler import --name=trusty-megamnode --path=/mnt --breed=ubuntu --arch x86_64

  echo "Imported trusty megam node iso to cobbler" >> $COBBLER_LOG

}

#--------------------------------------------------------------------------
#This function will print out boot menu
#--------------------------------------------------------------------------
setup_profile_trustynode() {
  cecho "Setting up profile trusty node..." $green
  echo  "Setting up profile trusty node..." >> $COBBLER_LOG

  cat > //etc/cobbler/pxe/pxedefault.template << 'EOF'
DEFAULT menu
PROMPT 0
MENU TITLE Megam Cloud In a Box(Node) | www.megam.io
TIMEOUT 200
TOTALTIMEOUT 6000
ONTIMEOUT $pxe_timeout_profile

$pxe_menu_items

MENU end
EOF

  wget -O /var/lib/cobbler/kickstarts/megamnode.seed http://get.megam.io/npreseed.cfg

  cp /usr/share/megam/megamcib/conf/trusty/cobblerd/install_post_cibnode.py /usr/lib/python2.7/dist-packages/cobbler/modules/
  cp /usr/share/megam/megamcib/conf/trusty/cobblerd/preseed_early_ub1404 /var/lib/cobbler/scripts/
  cp /usr/share/megam/megamcib/conf/trusty/cobblerd/preseed_late_ub1404 /var/lib/cobbler/scripts/
  cp /usr/share/megam/megamcib/conf/trusty/cobblerd/kickstart_start /var/lib/cobbler/snippets/
  cp /usr/share/megam/megamcib/conf/trusty/cobblerd/kickstart_done /var/lib/cobbler/snippets/
  cp /usr/share/megam/megamcib/conf/trusty/cobblerd/post_run_deb /var/lib/cobbler/snippets/

  cobbler profile edit --name="trusty-megamnode-x86_64" --kickstart="/var/lib/cobbler/kickstarts/megamnode.seed"
}

#--------------------------------------------------------------------------
#This function will add profile for ceph
#--------------------------------------------------------------------------
setup_profile_cephnode() {
  wget -O /var/lib/cobbler/kickstarts/megamceph.seed http://get.megam.io/cephpreseed.cfg
  cobbler profile add --name=trusty-megamceph --distro=trusty-megamnode-x86-64 --kickstart=/var/lib/cobbler/kickstarts/megamceph.seed
}


#--------------------------------------------------------------------------
#This function will print out an install report
#--------------------------------------------------------------------------
install_complete() {
  cecho "##################################################" $green
  cecho "cobbler installed successfully." $yellow
  cecho "Refer http://bit.ly/megamcib for more information." $yellow
  cecho "##################################################" $green
  echo  "Cobbler installed successfully." >> $COBBLER_LOG
}

#--------------------------------------------------------------------------
#This function will uninstall cobblerd
#--------------------------------------------------------------------------
uninstall_cobbler() {
  cecho "Uninstalling cobblerd.." $yellow

  apt-get -y remove cobbler cobbler-common cobbler-web  dnsmasq >> $COBBLER_LOG

  apt-get -y purge cobbler cobbler-common cobbler-web  >> $COBBLER_LOG

  apt-get -y remove debmirror >> $COBBLER_LOG

  apt-get -y purge debmirror >> $COBBLER_LOG

  apt-get -y remove xinetd tftpd tftp >> $COBBLER_LOG

  apt-get -y purge xinetd tftpd tftp >> $COBBLER_LOG

  [ -d /var/log/cobbler ] && rm -rf /var/log/cobbler
  [ -d /var/lib/cobbler ] && rm -rf /var/lib/cobbler
  [ -d /etc/cobbler ] && rm -rf /etc/cobbler

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
