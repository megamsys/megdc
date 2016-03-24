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

var debianverticeconfig *DebianVerticeConfig

func init() {
	debianverticeconfig = &DebianVerticeConfig{}
	templates.Register("DebianVerticeConfig", debianverticeconfig)
}

type DebianVerticeConfig struct{}

func (tpl *DebianVerticeConfig) Render(p urknall.Package) {
	p.AddTemplate("vertice", &DebianVerticeConfigTemplate{})
}

func (tpl *DebianVerticeConfig) Options(t *templates.Template) {
}

func (tpl *DebianVerticeConfig) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianVerticeConfig{})
}

type DebianVerticeConfigTemplate struct{}

func (m *DebianVerticeConfigTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("repository",
		u.Shell("echo 'deb [arch=amd64] "+DefaultMegamRepo+"' > "+ListFilePath),
		u.UpdatePackagesOmitError(),
	)

	pkg.AddCommands("vertice",
	u.InstallPackages("vertice"),
	)

}
