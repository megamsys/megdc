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

package debian

import (
	"github.com/megamsys/megdc/templates"
	u "github.com/megamsys/megdc/templates/ubuntu"
	"github.com/megamsys/urknall"
)

var debianonehostremove *DebianOneHostRemove

func init() {
	debianonehostremove = &DebianOneHostRemove{}
	templates.Register("DebianOneHostRemove", debianonehostremove)
}

type DebianOneHostRemove struct{}

func (tpl *DebianOneHostRemove) Render(p urknall.Package) {
	p.AddTemplate("onehost", &DebianOneHostRemoveTemplate{})
}

func (tpl *DebianOneHostRemove) Options(t *templates.Template) {
}

func (tpl *DebianOneHostRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianOneHostRemove{})
}

type DebianOneHostRemoveTemplate struct{}

func (m *DebianOneHostRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("onehost",
		u.RemovePackage("opennebula-node openvswitch-common openvswitch-switch bridge-utils sshpass"),
		u.RemovePackages(""),
		u.PurgePackages("opennebula-node openvswitch-common openvswitch-switch bridge-utils sshpass"),
	)
	pkg.AddCommands("onehost-clean",
		u.Shell("rm -r /var/lib/urknall/onehost*"),
	)
}
