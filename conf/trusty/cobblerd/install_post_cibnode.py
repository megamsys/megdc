#!/usr/bin/python
import distutils.sysconfig
import sys
import logging
import os
import traceback
import time

plib = distutils.sysconfig.get_python_lib()
mod_path="%s/cobbler" % plib
sys.path.insert(0, mod_path)

from utils import _
import smtplib
import sys
import cobbler.templar as templar
from cobbler.cexceptions import CX
import utils

logging.basicConfig(filename='/var/log/megam/megamcib/cobbler_triggers.log',level=logging.DEBUG)
logging.debug("%s\tRegistering node install trigger\n", time.strftime('%X %x %Z'))

def register():
   logging.debug("%s\tRegistered node install trigger\n", time.strftime('%X %x %Z'))
   # this pure python trigger acts as if it were a legacy shell-trigger, but is much faster.
   # the return of this method indicates the trigger type
   return "/var/lib/cobbler/triggers/install/post/*"

def run(api, args, logger):
    logging.debug("%s\tRunning node install trigger\n", time.strftime('%X %x %Z'))
    settings = api.settings()

    # go no further if this feature is turned off
    if not str(settings.base_megamreporting_enabled).lower() in [ "1", "yes", "y", "true"]:
        return 0

    logging.debug("%s\tUnwrapping args\n", time.strftime('%X %x %Z'))

    objtype = args[0] # "target" or "profile"
    name    = args[1] # name of target or profile
    box_ip = args[2] # ip

    logging.debug("%s\tArgs are \t%s\t%s\t%s\n", time.strftime('%X %x %Z'), objtype, name, box_ip)

    if objtype == "system":
        target = api.find_system(name)
    else:
        target = api.find_profile(name)

    logging.debug("%s\tFigured out objtype\t%s\n", time.strftime('%X %x %Z'),target)

    target = utils.blender(api, False, target)

    if target == {}:
        raise CX("failure looking up target")


    logging.debug("%s\tboxip is\t%s\n", time.strftime('%X %x %Z'),box_ip)


    with open('/var/lib/megam/megamcib/boxips', 'a') as f:
      	f.write(box_ip)


    logging.debug("%s\tI am done. Adios..Amigo\n",time.strftime('%X %x %Z'))

    return 0
