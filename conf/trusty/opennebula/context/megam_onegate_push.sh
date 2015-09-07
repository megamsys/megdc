#!/bin/bash
#Copyright (c) 2014 Megam Systems.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
###############################################################################
# The script pushes the public ipaddress using the command
# ifconfig ham0 | grep inet addr and sends it back in the variable IP_ADDRESS
# change the ham0 to your specific network interface (eth1, eth2, eth0 etc.)
# When TOKEN=YES is set in the template in opennebula automatically populates
# $ONEGATE_TOKEN and $ONEGATE_URL.
# Note: this script gets called from init.sh
###############################################################################

ERROR=0

if [ -z $ONEGATE_TOKEN ]; then
    echo "ONEGATE_TOKEN env variable must point to the token.txt file"
    ERROR=1
fi

if [ -z $ONEGATE_URL ]; then
    echo "ONEGATE_URL env variable must be set"
    ERROR=1
fi

if [ $ERROR = 1 ]; then
    exit -1
fi

curl -X "PUT" --header "X-ONEGATE-TOKEN: `cat $ONEGATE_TOKEN`" $ONEGATE_URL -d "IP_ADDRESS=`ifconfig ham0 | grep "inet addr" | awk -F: '{print $2}' | awk '{print $1}'`"
