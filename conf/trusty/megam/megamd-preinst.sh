#!/bin/bash

#Gem install
gem install chef --no-ri --no-rdoc
mkdir -p /var/lib/megam/gems
cd /var/lib/megam/gems

: <<'END'
set -- https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-ec2.gem https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-gogrid-0.1.0.gem https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-google-1.3.1.gem https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-hp.gem https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-opennebula-0.2.0.gem https://s3-ap-southeast-1.amazonaws.com/megampub/gems/profitbricks-1.1.0.20130930123843.gem https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-profitbricks-0.3.0.gem
while [ $# -gt 0 ]
do
        wget $1
        gem install ${1##*/} --no-ri --no-rdoc
        shift;
done
END

wget https://s3-ap-southeast-1.amazonaws.com/megampub/gems/knife-opennebula-0.2.0.gem

gem install knife-opennebula-0.2.0.gem

if [ -d "/opt/chef-server" ]; then

sudo chef-server-ctl reconfigure

#Rabbitmq server has to be run in localhost
#cat /etc/rabbitmq/rabbitmq-env.conf	NODENAME=localhost

if [ -d "/etc/rabbitmq" ]; then

#Rabbitmq server has to be run in localhost for chef-server changes
#cat /etc/rabbitmq/rabbitmq-env.conf	NODENAME=localhost

cat > //etc/rabbitmq/rabbitmq-env.conf <<EOF
NODENAME=localhost
EOF
sudo rabbitmqctl stop_app
sudo rabbitmqctl reset
sudo rabbitmqctl stop

fi

getip(){
while read Iface Destination Gateway Flags RefCnt Use Metric Mask MTU Window IRTT; do
		[ "$Mask" = "00000000" ] && \
		interface="$Iface" && \
		ipaddr=$(LC_ALL=C /sbin/ip -4 addr list dev "$interface" scope global) && \
		ipaddr=${ipaddr#* inet } && \
		ipaddr=${ipaddr%%/*} && \
		break
	done < /proc/net/route
}
getip

cat > //etc/chef-server/chef-server.rb <<EOF
nginx['url']="https://$ipaddr"
lb['api_fqdn']="$ipaddr"
lb['web_ui_fqdn']="$ipaddr"
nginx['server_name']="$ipaddr"
nginx['non_ssl_port'] = 90

EOF

sudo chef-server-ctl reconfigure

sudo chef-server-ctl restart

sudo rabbitmq-server -detached

  set -e

  #chef_repo_dir=`find /var/lib/megam/megamd  -name chef-repo  | awk -F/ -vOFS=/ 'NF-=0' | sort -u`
   chef_repo_dir="/var/lib/megam/megamd/"

git clone https://github.com/megamsys/chef-repo.git $chef_repo_dir
  cp /etc/chef-server/admin.pem $chef_repo_dir/chef-repo/.chef
  cp /etc/chef-server/chef-validator.pem $chef_repo_dir/chef-repo/.chef
  
 sed -i "s@^[ \t]*chef_server_url.*@chef_server_url 'https://$ipaddr'@" $chef_repo_dir/chef-repo/.chef/knife.rb
  
mkdir $chef_repo_dir/chef-repo/.chef/trusted_certs
cp /var/opt/chef-server/nginx/ca/$ipaddr.crt /var/lib/megam/megamd/chef-repo/.chef/trusted_certs

#chown -R cibadmin:cibadmin $chef_repo_dir/chef-repo
  knife cookbook upload --all -c $chef_repo_dir/chef-repo/.chef

fi
