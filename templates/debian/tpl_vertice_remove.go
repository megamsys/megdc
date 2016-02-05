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

var debianverticeremove *DebianMegamdRemove

func init() {
	debianverticeremove = &DebianMegamdRemove{}
	templates.Register("DebianMegamdRemove", debianverticeremove)
}

type DebianMegamdRemove struct{}

func (tpl *DebianMegamdRemove) Render(p urknall.Package) {
	p.AddTemplate("vertice", &DebianMegamdRemoveTemplate{})
}

func (tpl *DebianMegamdRemove) Options(t *templates.Template) {
}

func (tpl *DebianMegamdRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianMegamdRemove{})
}

type DebianMegamdRemoveTemplate struct{}

func (m *DebianMegamdRemoveTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("vertice",
		u.RemovePackage("vertice"),
		u.RemovePackages(""),
		u.PurgePackages("vertice"),
		u.Shell("dpkg --get-selections megam*"),
	)
	pkg.AddCommands("vertice-clean",
		u.Shell("rm -r /var/lib/urknall/vertice*"),
	)
}
