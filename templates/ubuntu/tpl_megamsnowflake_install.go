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

var ubuntusnowflakeinstall *UbuntuSnowflakeInstall

func init() {
	ubuntusnowflakeinstall = &UbuntuSnowflakeInstall{}
	templates.Register("UbuntuSnowflakeInstall", ubuntusnowflakeinstall)
}

type UbuntuSnowflakeInstall struct{}

func (tpl *UbuntuSnowflakeInstall) Render(p urknall.Package) {
	p.AddTemplate("snowflake", &UbuntuSnowflakeInstallTemplate{})
}

func (tpl *UbuntuSnowflakeInstall) Options(t *templates.Template) {
}

func (tpl *UbuntuSnowflakeInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuSnowflakeInstall{})
}

type UbuntuSnowflakeInstallTemplate struct{}

func (m *UbuntuSnowflakeInstallTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("repository",
		Shell("echo 'deb [arch=amd64] "+DefaultMegamRepo+"' > "+ListFilePath),
		UpdatePackagesOmitError(),
	)

	pkg.AddCommands("megamsnowflake",
		InstallPackages("megamsnowflake"),
	)
}
