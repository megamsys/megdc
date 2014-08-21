#!/bin/bash

echo "--------------------megam install----------------"


apt-get -y install megamnilavu >> /var/log/megam.log

apt-get -y install megamgateway >> /var/log/megam.log

apt-get -y install megamd >> /var/log/megam.log

apt-get -y install megamanalytics >> /var/log/megam.log


