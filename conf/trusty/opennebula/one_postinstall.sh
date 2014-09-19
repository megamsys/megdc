#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"
echo "One POSTINSTALL install gems=======> " >> $ONE_INSTALL_LOG

#sed -i 's/.*dependencies_common => .*/:dependencies_common => [],/' /usr/share/one/install_gems
sudo rm /usr/share/one/install_gems

#/var/lib/megam/megamcib/install_gems has to be created early
sudo cp /usr/share/megam/megamcib/conf/trusty/opennebula/install_gems /usr/share/one/install_gems

sudo chmod 755 /usr/share/one/install_gems

sudo /usr/share/one/install_gems sunstone >> $ONE_INSTALL_LOG

