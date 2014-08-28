#!/bin/bash

echo "--------------------megam install----------------"
#Testing for error
pwd >> /var/log/megam/megam.log
ls >> /var/log/megam/megam.log
pwd >> /var/log/megam/megam.log
exit 5 >> /var/log/megam/megam.log

apt-get -y install megamnilavu >> /var/log/megam/megam.log


apt-get -y install megamgateway >> /var/log/megam/megam.log

apt-get -y install megamd >> /var/log/megam/megam.log

apt-get -y install megamanalytics >> /var/log/megam/megam.log


