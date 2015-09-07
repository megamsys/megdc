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
# The scripts sets up a VPN if NEW_VPN is passed in the context.
# This assumes that we use logmein-hamachi vpn : https://secure.logmein.com/products/hamachi/
# and can be downloaded from
# https://secure.logmein.com/labs/logmein-hamachi_2.1.0.119-1_amd64.deb
# The image we are running has logmein packaged inside.
# a) When NEW_VPN=YES: A new vpn network using VPN_NAME and VPN_PASSWORD is created.
# b) When NEW_VPN=NO : We join into the vpn network  VPN_NAME
# c) This script is ignored if NEW_VPN isn't passed in the CONTEXT
# Note : This script is called from init.sh
###############################################################################
if [ -z "$NEW_VPN" ]
  then
  echo "Found NEW_VPN"
  if [[ $NEW_VPN = "YES" ]]; then
    hamachi login
    hamachi create $VPN_NAME $VPN_PASSWORD
  else
    hamachi login
    hamachi do-join $VPN_NAME $VPN_PASSWORD
  fi
fi
