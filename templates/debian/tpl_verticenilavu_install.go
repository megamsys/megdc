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

const (
	// DefaultMegamRepo is the default megam repository to install if its not provided.
	DefaultMegamRepo = "http://get.megam.io/0.9/debian/8/ jessie testing"

	ListFilePath = "/etc/apt/sources.list.d/megam.list"
)

var debiannilavuinstall *DebianNilavuInstall

func init() {
	debiannilavuinstall = &DebianNilavuInstall{}
	templates.Register("DebianNilavuInstall", debiannilavuinstall)
}

type DebianNilavuInstall struct{}

func (tpl *DebianNilavuInstall) Render(p urknall.Package) {
	p.AddTemplate("nilavu", &DebianNilavuInstallTemplate{})
}

func (tpl *DebianNilavuInstall) Options(t *templates.Template) {
}

func (tpl *DebianNilavuInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianNilavuInstall{})
}

type DebianNilavuInstallTemplate struct{}

func (m *DebianNilavuInstallTemplate) Render(pkg urknall.Package) {
  //fail on ruby2.0 < check

	pkg.AddCommands("repository",
		u.Shell("echo 'deb [arch=amd64] "+DefaultMegamRepo+"' > "+ListFilePath),
		u.UpdatePackagesOmitError(),
	)

	pkg.AddCommands("verticenilavu",
		u.InstallPackages(" sdfaskld verticenilavu"),
	)

}
