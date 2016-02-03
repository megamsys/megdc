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

var debianoneremove *DebianOneRemove

func init() {
	debianoneremove = &DebianOneRemove{}
	templates.Register("DebianOneRemove", debianoneremove)
}

type DebianOneRemove struct{}

func (tpl *DebianOneRemove) Render(p urknall.Package) {
	p.AddTemplate("one", &DebianOneRemoveTemplate{})
}

func (tpl *DebianOneRemove) Options(t *templates.Template) {
}

func (tpl *DebianOneRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianOneRemove{})
}

type DebianOneRemoveTemplate struct{}

func (m *DebianOneRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("one",
		u.RemovePackage("opennebula opennebula-sunstone"),
		u.RemovePackages(""),
		u.PurgePackages("opennebula opennebula-sunstone"),
	)
	pkg.AddCommands("one-clean",
		u.Shell("rm -r /var/lib/urknall/one.*"),
	)
}
