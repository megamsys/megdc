#!/bin/bash

if [ -f /etc/lsb-release ]; then
    . /etc/lsb-release
    DISTRO=$DISTRIB_ID
elif [ -f /etc/debian_version ]; then
    DISTRO=Debian
    # XXX or Ubuntu
elif [ -f /etc/redhat-release ]; then
    DISTRO="Red Hat"
    # XXX or CentOS or Fedora
else
    DISTRO=$(uname -s)
fi



if [ "$DISTRO" = "Red Hat" ]  || [ "$DISTRO" = "Ubuntu" ] || [ "$DISTRO" = "Debian" ]
then
CONF='//usr/share/megam/verticegulpd/conf/gulpd.conf'
else if [ "$DISTRO" = "CoreOS" ]; then
  CONF='//var/lib/megam/verticegulpd/conf/gulpd.conf'
fi
fi
cat >$CONF  << 'EOF'

### Welcome to the Gulpd configuration file.

  ###
  ### [meta]
  ###
  ### Controls the parameters for the Raft consensus group that stores metadata
  ### about the gulp.
  ###

  [meta]
  [meta]
      user = "root"
      nsqd = ["192.168.1.100:4150"]
      scylla = ["192.168.1.100"]
      scylla_keyspace = "vertice"

  ###
  ### [gulpd]
  ###
  ### Controls which assembly to be deployed into machine
  ###

  [gulpd]
    enabled =true
    name_gulp = "hostname"
    cats_id = "AMS1259077729232486400"
    cat_id = "ASM1260230009767985152"
	provider = "chefsolo"
	cookbook = "megam_run"
	repository = "github"
	repository_path = "https://github.com/megamsys/chef-repo.git"
  repository_tar_path = "https://github.com/megamsys/chef-repo/archive/0.96.tar.gz"

  ###
  ### [http]
  ###
  ### Controls how the HTTP endpoints are configured. This a frill
  ### mechanism for pinging gulpd (ping)
  ###

  [http]
    enabled = false
    bind-address = "127.0.0.1:6666"

EOF

sed -i "s/^[ \t]*name_gulp.*/    name = \"$NODE_NAME\"/" $CONF
sed -i "s/^[ \t]*cats_id.*/    cats_id = \"$ASSEMBLIES_ID\"/" $CONF
sed -i "s/^[ \t]*cat_id.*/    cat_id = \"$ASSEMBLY_ID\"/" $CONF

case "$DISTRO" in
   "Ubuntu")
stop verticegulpd
start verticegulpd
   ;;
   "Debian")
systemctl stop verticegulpd.service
systemctl start verticegulpd.service
systemctl stop cadvisor.service
systemctl start cadvisor.service
   ;;
   "Red Hat")
systemctl stop verticegulpd.service
systemctl start verticegulpd.service
systemctl stop cadvisor.service
systemctl start cadvisor.service
   ;;
   "CoreOS")
if [ -f /mnt/context.sh ]; then
  . /mnt/context.sh
fi

sudo cat >> //home/core/.ssh/authorized_keys <<EOF
$SSH_PUBLIC_KEY
EOF

sudo -s

sudo cat > //etc/hostname <<EOF
$HOSTNAME
EOF

sudo cat > //etc/hosts <<EOF
$ETHO_IP $HOSTNAME localhost

EOF
sudo cat > //etc/systemd/network/static.network <<EOF
[Match]
Name=ens3

[Network]
Address=$ETH0_IP/24
Gateway=$ETH0_GATEWAY
DNS=8.8.8.8
DNS=8.8.4.4
EOF

sudo systemctl restart systemd-networkd
systemctl stop verticegulpd.service
systemctl start verticegulpd.service
systemctl stop cadvisor.service
systemctl start cadvisor.service


   ;;
esac
