#!/bin/bash
#root@megammaster:/onemegam/coreos# cat init.sh 
#root@megammaster:/megam# cat init.sh 
 
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

sudo cat >> //etc/hosts <<EOF
$IP_ADDRESS $HOSTNAME localhost

EOF

sudo cat > //etc/systemd/network/static.network <<EOF
[Match]
Name=ens3

[Network]
Address=$IP_ADDRESS/24
Gateway=$GATEWAY
EOF

sudo systemctl restart systemd-networkd
