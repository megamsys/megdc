#!/bin/bash

ONE_INSTALL_LOG="/var/log/megam/megamcib/one_install.log"

/usr/share/one/install_gems sunstone >> $ONE_INSTALL_LOG
