/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package ubuntu

import (
	"github.com/dynport/urknall"
	"github.com/megamsys/megdc/templates"
)
<<<<<<< HEAD
const{
 Host = "hostname"
}
=======

>>>>>>> origin/master
var ubuntucephremove *UbuntuCephRemove

func init() {
	ubuntucephremove = &UbuntuCephRemove{}
	templates.Register("UbuntuCephRemove", ubuntucephremove)
}

type UbuntuCephRemove struct{}

func (tpl *UbuntuCephRemove) Render(p urknall.Package) {
	p.AddTemplate("nilavu", &UbuntuCephRemoveTemplate{})
}

func (tpl *UbuntuCephRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCephRemove{})
}

type UbuntuCephRemoveTemplate struct{}

func (m *UbuntuCephRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("purgedata",
<<<<<<< HEAD
		Shell("ceph-deploy purgedata `" + Host + "`"),
=======
		Shell("ceph-deploy purgedata " + Host + ""),
>>>>>>> origin/master
	)
  pkg.AddCommands("forgetKeys",
		Shell("ceph-deploy forgetkeys"),
	)
  pkg.AddCommands("purge",
		Shell("ceph-deploy purge " + Host + ""),
	)
  pkg.AddCommands("remove",
		Shell("sudo rm -r /var/lib/ceph/"),
	)
  pkg.AddCommands("remove",
		Shell("sudo rm -r /var/lib/ceph/"),
	)
}
