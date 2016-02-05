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

var debiangatewayinstall *DebianGatewayInstall

func init() {
	debiangatewayinstall = &DebianGatewayInstall{}
	templates.Register("DebianGatewayInstall", debiangatewayinstall)
}

type DebianGatewayInstall struct{}

func (tpl *DebianGatewayInstall) Render(p urknall.Package) {
	p.AddTemplate("gateway", &DebianGatewayInstallTemplate{})
}

func (tpl *DebianGatewayInstall) Options(t *templates.Template) {
}

func (tpl *DebianGatewayInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianGatewayInstall{})
}

type DebianGatewayInstallTemplate struct{}

func (m *DebianGatewayInstallTemplate) Render(pkg urknall.Package) {
	//fail on Java -version (1.8 check)
	pkg.AddCommands("repository",
		u.Shell("echo 'deb [arch=amd64] "+DefaultMegamRepo+"' > "+ListFilePath),
		u.UpdatePackagesOmitError(),
	)

	pkg.AddCommands("verticegateway",
		u.InstallPackages("verticegateway"),
	)
}
