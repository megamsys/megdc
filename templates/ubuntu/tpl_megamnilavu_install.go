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

var ubuntumegamnilavuinstall *UbuntuMegamNilavuInstall

func init() {
	ubuntumegamnilavuinstall = &UbuntuMegamNilavuInstall{}
	templates.Register("UbuntuMegamNilavuInstall", ubuntumegamnilavuinstall)
}

type UbuntuMegamNilavuInstall struct{}

func (tpl *UbuntuMegamNilavuInstall) Render(p urknall.Package) {
	p.AddTemplate("nilavu", &UbuntuMegamNilavuInstallTemplate{})
}

func (tpl *UbuntuMegamNilavuInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamNilavuInstall{})
}

type UbuntuMegamNilavuInstallTemplate struct{}

func (m *UbuntuMegamNilavuInstallTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("packages",
		InstallPackages("software-properties-common", "python-software-properties", "ruby2.0", "ruby2.0-dev"),
	)
	
	pkg.AddCommands("repository",
		Shell("add-apt-repository 'deb [arch=amd64] http://get.megam.io/0.9/ubuntu/14.04/ testing megam'"),
		Shell("apt-key adv --keyserver keyserver.ubuntu.com --recv B3E0C1B7"),
		UpdatePackagesOmitError(),
	)
	
	pkg.AddCommands("megamcommon",
		And("apt-get -y install megamcommon"),
	)

	pkg.AddCommands("megamnilavu",
		InstallPackages("megamnilavu"),
	)
}
