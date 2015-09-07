#!/bin/bash

#root@megammaster:/megam# cat init.sh 
#stop megamgulpd
#mkdir /$NODE_NAME


cat > //usr/share/megam/megamgulpd/conf/gulpd.conf << 'EOF'
megam_home: /var/lib/megam/
account_id: ACT1242476978897027072
name: $NODE_NAME
id: ASM1242477979045593088
docker_path: /var/lib/docker/containers/
riak:
  url: 192.168.1.100:8087
  bucket: catreqs
api:
  host: api.megam.io
amqp:
  url: amqp://guest:guest@192.168.1.105:5672/
  exchange: megam_reapportions.megam.co_exchange
  queue: megam_reapportions.megam.co_queue
  consumerTag: megam_node_consumer
  routingkey: megam_key
admin:
  port: 8084
EOF
stop megamgulpd
start megamgulpd

sudo echo 3 > /proc/sys/vm/drop_caches
