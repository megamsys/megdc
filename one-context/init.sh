#!/bin/bash
#root@megammaster:/megam# cat init.sh
#stop megamgulpd
#mkdir /$NODE_NAME
#mkdir /$ASSEMBLY_ID

cat > //usr/share/megam/megamgulpd/conf/gulpd.conf << 'EOF'

### Welcome to the Gulpd configuration file.

  ###
  ### [meta]
  ###
  ### Controls the parameters for the Raft consensus group that stores metadata
  ### about the gulp.
  ###

  [meta]
    debug = true
    hostname = "localhost"
    bind_address = "192.168.1.105:7777"
    dir = "/var/lib/megam/gulp/meta"
    riak = ["192.168.1.105:8087"]
    api  = "https://api.megam.io/v2"
    amqp = "amqp://guest:guest@192.168.1.105:5672/"

  ###
  ### [gulpd]
  ###
  ### Controls which assembly to be deployed into machine
  ###

  [gulpd]
    name_gulp = "hostname"
    cats_id = "AMS1259077729232486400"
    cat_id = "ASM1260230009767985152"
	provider = "chefsolo"
	cookbook = "megam_run"
	repository = "github"
	repository_path = "https://github.com/megamsys/chef-repo.git"

  ###
  ### [http]
  ###
  ### Controls how the HTTP endpoints are configured. This a frill
  ### mechanism for pinging gulpd (ping)
  ###

  [http]
    enabled = true
    bind-address = "localhost:6666"

EOF

sed -i "s/^[ \t]*name_gulp.*/    name = \"$NODE_NAME\"/" /usr/share/megam/megamgulpd/conf/gulpd.conf
sed -i "s/^[ \t]*cats_id.*/    cats_id = \"$ASSEMBLIES_ID\"/" /usr/share/megam/megamgulpd/conf/gulpd.conf
sed -i "s/^[ \t]*cat_id.*/    cat_id = \"$ASSEMBLY_ID\"/" /usr/share/megam/megamgulpd/conf/gulpd.conf



stop megamgulpd
start megamgulpd

sudo echo 3 > /proc/sys/vm/drop_caches
