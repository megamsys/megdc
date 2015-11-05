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

const (
	// DefaultMegamRepo is the default megam repository to install if its not provided.
	    DefaultMegamRepo = "http://get.megam.io/0.9/ubuntu/14.04/ trusty testing"

			ListFilePath  = "/etc/apt/sources.list.d/megam.list"
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
		pkg.AddCommands("repository",
		Shell("echo 'deb [arch=amd64] " + DefaultMegamRepo + "' > " + ListFilePath ),
		UpdatePackagesOmitError(),
	)

	pkg.AddCommands("megamcommon",
		And("apt-get -y install megamcommon"),
	)

	pkg.AddCommands("megamnilavu",
		InstallPackages("megamnilavu"),
	)
}
