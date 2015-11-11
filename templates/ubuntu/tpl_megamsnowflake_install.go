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
	"github.com/megamsys/urknall"
	"github.com/megamsys/megdc/templates"
)


var ubuntumegamsnowflakeinstall *UbuntuMegamSnowflakeInstall

func init() {
	ubuntumegamsnowflakeinstall = &UbuntuMegamSnowflakeInstall{}
	templates.Register("UbuntuMegamSnowflakeInstall", ubuntumegamsnowflakeinstall)
}

type UbuntuMegamSnowflakeInstall struct{}

func (tpl *UbuntuMegamSnowflakeInstall) Render(p urknall.Package) {
	p.AddTemplate("snowflake", &UbuntuMegamSnowflakeInstallTemplate{})
}

func (tpl *UbuntuMegamSnowflakeInstall) Options(opts map[string]string) {
}

func (tpl *UbuntuMegamSnowflakeInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamSnowflakeInstall{})
}

type UbuntuMegamSnowflakeInstallTemplate struct{}

func (m *UbuntuMegamSnowflakeInstallTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("repository",
		Shell("echo 'deb [arch=amd64] " + DefaultMegamRepo + "' > " + ListFilePath),
		UpdatePackagesOmitError(),
	)

	pkg.AddCommands("megamsnowflake",
		InstallPackages("megamsnowflake"),
	)
}
