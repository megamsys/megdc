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

var debianmegamdremove *DebianMegamdRemove

func init() {
	debianmegamdremove = &DebianMegamdRemove{}
	templates.Register("DebianMegamdRemove", debianmegamdremove)
}

type DebianMegamdRemove struct{}

func (tpl *DebianMegamdRemove) Render(p urknall.Package) {
	p.AddTemplate("megamd", &DebianMegamdRemoveTemplate{})
}

func (tpl *DebianMegamdRemove) Options(t *templates.Template) {
}

func (tpl *DebianMegamdRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianMegamdRemove{})
}

type DebianMegamdRemoveTemplate struct{}

func (m *DebianMegamdRemoveTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("megamd",
		u.RemovePackage("megamd"),
		u.RemovePackages(""),
		u.PurgePackages("megamd"),
		u.Shell("dpkg --get-selections megam*"),
	)
	pkg.AddCommands("megamd-clean",
		u.Shell("rm -r /var/lib/urknall/megamd*"),
	)
}
