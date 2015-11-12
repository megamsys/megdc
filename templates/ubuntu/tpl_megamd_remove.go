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
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"
)

var ubuntumegamdremove *UbuntuMegamdRemove

func init() {
	ubuntumegamdremove = &UbuntuMegamdRemove{}
	templates.Register("UbuntuMegamdRemove", ubuntumegamdremove)
}

type UbuntuMegamdRemove struct{}

func (tpl *UbuntuMegamdRemove) Render(p urknall.Package) {
	p.AddTemplate("megamd", &UbuntuMegamdRemoveTemplate{})
}

func (tpl *UbuntuMegamdRemove) Options(opts map[string]string) {
}

func (tpl *UbuntuMegamdRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamdRemove{})
}

type UbuntuMegamdRemoveTemplate struct{}

func (m *UbuntuMegamdRemoveTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("megamd",
		RemovePackage("megamd"),
		RemovePackages(""),
		PurgePackages("megamd"),
		Shell("dpkg --get-selections megam*"),
	)
}
