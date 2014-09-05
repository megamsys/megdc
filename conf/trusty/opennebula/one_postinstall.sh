#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/opennebula.log"

/usr/share/one/install_gems sunstone >> $ONE_INSTALL_LOG
