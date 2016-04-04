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

package centos

import (
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"

)

var centosonehostremove *CentOsOneHostRemove

func init() {
	centosonehostremove = &CentOsOneHostRemove{}
	templates.Register("CentOsOneHostRemove", centosonehostremove)
}

type CentOsOneHostRemove struct{}

func (tpl *CentOsOneHostRemove) Render(p urknall.Package) {
	p.AddTemplate("onehost", &CentOsOneHostRemoveTemplate{})
}

func (tpl *CentOsOneHostRemove) Options(t *templates.Template) {
}

func (tpl *CentOsOneHostRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &CentOsOneHostRemove{})
}

type CentOsOneHostRemoveTemplate struct{}

func (m *CentOsOneHostRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("onehost",
		RemovePackage("opennebula-node openvswitch-common openvswitch-switch bridge-utils sshpass"),
		RemovePackages(""),
		PurgePackages("opennebula-node openvswitch-common openvswitch-switch bridge-utils sshpass"),
	)
	pkg.AddCommands("onehost-clean",
		Shell("rm -rf /var/lib/urknall/onehost*"),
	)
}
