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

var debiannilavuremove *DebianNilavuRemove

func init() {
	debiannilavuremove = &DebianNilavuRemove{}
	templates.Register("DebianNilavuRemove", debiannilavuremove)
}

type DebianNilavuRemove struct{}

func (tpl *DebianNilavuRemove) Render(p urknall.Package) {
	p.AddTemplate("nilavu", &DebianNilavuRemoveTemplate{})
}

func (tpl *DebianNilavuRemove) Options(t *templates.Template) {
}

func (tpl *DebianNilavuRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianNilavuRemove{})
}

type DebianNilavuRemoveTemplate struct{}

func (m *DebianNilavuRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("verticenilavu",
		u.RemovePackage("verticenilavu"),
		u.RemovePackages(""),
		u.PurgePackages("verticenilavu"),
		u.Shell("dpkg --get-selections megam*"),
	)
	pkg.AddCommands("nilavu-clean",
		u.Shell("rm -r /var/lib/urknall/nilavu*"),
	)
}
