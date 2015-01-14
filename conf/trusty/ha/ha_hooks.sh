#!/bin/sh

cat <<EOT >> /etc/one/oned.conf 
VM_HOOK = [
   name      = "on_failure_recreate",
   on        = "STOP",
   command   = "/usr/bin/env onevm delete --recreate",
   arguments = "$ID $TEMPLATE" ]


HOST_HOOK = [
    name      = "host_error",
    on        = "DISABLE",
    command   = "/var/lib/one/remotes/hooks/ft/host_error.rb",
    arguments = "$ID -r",
    remote    = "no" ]
EOT


sunstone-server restart
econe-server restart
occi-server restart
onegate-server restart
one restart
sudo service opennebula restart
