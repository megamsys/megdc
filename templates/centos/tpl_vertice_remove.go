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

var centosverticeremove *CentOsMegamdRemove

func init() {
	centosverticeremove = &CentOsMegamdRemove{}
	templates.Register("CentosMegamdRemove", centosverticeremove)
}

type CentOsMegamdRemove struct{}

func (tpl *CentOsMegamdRemove) Render(p urknall.Package) {
	p.AddTemplate("vertice", &CentOsMegamdRemoveTemplate{})
}

func (tpl *CentOsMegamdRemove) Options(t *templates.Template) {
}

func (tpl *CentOsMegamdRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &CentOsMegamdRemove{})
}

type CentOsMegamdRemoveTemplate struct{}

func (m *CentOsMegamdRemoveTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("vertice",
		RemovePackage("vertice"),
		RemovePackages(""),
		PurgePackages("vertice"),
		Shell("dpkg --get-selections megam*"),
	)
	pkg.AddCommands("vertice-clean",
		Shell("rm -rf /var/lib/urknall/vertice*"),
	)
}
