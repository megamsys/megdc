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

const (
	// DefaultMegamRepo is the default megam repository to install if its not provided.
	DefaultMegamRepo = "http://get.megam.io/0.9/ubuntu/14.04/ trusty nightly"

	ListFilePath = "/etc/apt/sources.list.d/megam.list"
)

var ubuntunilavuinstall *UbuntuNilavuInstall

func init() {
	ubuntunilavuinstall = &UbuntuNilavuInstall{}
	templates.Register("UbuntuNilavuInstall", ubuntunilavuinstall)
}

type UbuntuNilavuInstall struct{}

func (tpl *UbuntuNilavuInstall) Render(p urknall.Package) {
	p.AddTemplate("nilavu", &UbuntuNilavuInstallTemplate{})
}

func (tpl *UbuntuNilavuInstall) Options(t *templates.Template) {
}

func (tpl *UbuntuNilavuInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuNilavuInstall{})
}

type UbuntuNilavuInstallTemplate struct{}

func (m *UbuntuNilavuInstallTemplate) Render(pkg urknall.Package) {
  //fail on ruby2.0 < check

	pkg.AddCommands("repository",
		Shell("echo 'deb [arch=amd64] "+DefaultMegamRepo+"' > "+ListFilePath),
		UpdatePackagesOmitError(),
	)

	pkg.AddCommands("verticenilavu",
		InstallPackages("verticenilavu"),
	)

}
