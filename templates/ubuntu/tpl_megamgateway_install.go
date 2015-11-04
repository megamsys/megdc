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

var ubuntumegamgatewayinstall *UbuntuMegamGatewayInstall

func init() {
	ubuntumegamgatewayinstall = &UbuntuMegamGatewayInstall{}
	templates.Register("UbuntuMegamGatewayInstall", ubuntumegamgatewayinstall)
}

type UbuntuMegamGatewayInstall struct{}

func (tpl *UbuntuMegamGatewayInstall) Render(p urknall.Package) {
	p.AddTemplate("gateway", &UbuntuMegamGatewayInstallTemplate{})
}

func (tpl *UbuntuMegamGatewayInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamGatewayInstall{})
}

type UbuntuMegamGatewayInstallTemplate struct{}

func (m *UbuntuMegamGatewayInstallTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("repository",
		Shell(" echo 'deb [arch=amd64] http://get.megam.io/0.9/ubuntu/14.04/ trusty testing' > /etc/apt/sources.list.d/megam.list"),
		UpdatePackagesOmitError(),
	)

	pkg.AddCommands("megamgateway",
		InstallPackages("megamgateway"),
	)
}
