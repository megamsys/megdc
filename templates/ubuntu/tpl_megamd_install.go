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

var ubuntumegamdinstall *UbuntuMegamdInstall

func init() {
	ubuntumegamdinstall = &UbuntuMegamdInstall{}
	templates.Register("UbuntuMegamdInstall", ubuntumegamdinstall)
}

type UbuntuMegamdInstall struct{}

func (tpl *UbuntuMegamdInstall) Render(p urknall.Package) {
	p.AddTemplate("megamd", &UbuntuMegamdInstallTemplate{})
}

func (tpl *UbuntuMegamdInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamdInstall{})
}

type UbuntuMegamdInstallTemplate struct{}

func (m *UbuntuMegamdInstallTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("build",
		And("ls -la",),
	)
}
