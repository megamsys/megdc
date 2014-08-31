#!/usr/bin/python
import distutils.sysconfig
import sys
import os
import traceback

plib = distutils.sysconfig.get_python_lib()
mod_path="%s/cobbler" % plib
sys.path.insert(0, mod_path)

from utils import _
import smtplib
import sys
import cobbler.templar as templar
from cobbler.cexceptions import CX
import utils


def register():
   # this pure python trigger acts as if it were a legacy shell-trigger, but is much faster.
   # the return of this method indicates the trigger type
   return "/var/lib/cobbler/triggers/install/post/*"

def run(api, args, logger):
    
    settings = api.settings()

    # go no further if this feature is turned off
    if not str(settings.base_megamreporting_enabled).lower() in [ "1", "yes", "y", "true"]:
        return 0

    objtype = args[0] # "target" or "profile"
    name    = args[1] # name of target or profile
    boot_ip = args[2] # ip or "?"

    if logger is not None:
  	logger.warning("post install for cib node [objtype=%s, name=%s, boot_ip=%s]",objtype, name, boot_ip)


    if objtype == "system":
        target = api.find_system(name)
    else:
        target = api.find_profile(name)

    if logger is not None:
	logger.warning("[post install] %s -> target is %s", name, target)

    # collapse the object down to a rendered datastructure
    target = utils.blender(api, False, target)

    if target == {}:
        raise CX("failure looking up target")

    boxip = target.get_ip_address()

    if logger is not None:
	logger.warning("[post install] %s -> boxip is %s", target, boxip)

    with open('/var/log/megam/megamcib/boxips', 'a') as f:
	f.write(name + "=" + boxip)
    
    if logger is not None:
	logger.warning("[post install] %s wrote \"%s\" to boxips file", target, boxip)
    
    return 0




