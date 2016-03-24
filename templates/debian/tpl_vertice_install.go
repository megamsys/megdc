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

var debianverticeinstall *DebianMegamdInstall

func init() {
	debianverticeinstall = &DebianMegamdInstall{}
	templates.Register("DebianMegamdInstall", debianverticeinstall)
}

type DebianMegamdInstall struct{}

func (tpl *DebianMegamdInstall) Render(p urknall.Package) {
	p.AddTemplate("vertice", &DebianMegamdInstallTemplate{})
}

func (tpl *DebianMegamdInstall) Options(t *templates.Template) {
}

func (tpl *DebianMegamdInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianMegamdInstall{})
}

type DebianMegamdInstallTemplate struct{}

func (m *DebianMegamdInstallTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("repository",
		u.Shell("echo 'deb [arch=amd64] "+DefaultMegamRepo+"' > "+ListFilePath),
		u.UpdatePackagesOmitError(),
	)

	pkg.AddCommands("vertice",
	u.InstallPackages("vertice"),
	)

}
