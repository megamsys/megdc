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

var ubuntumegamsnowflakeremove *UbuntuMegamSnowflakeRemove

func init() {
	ubuntumegamsnowflakeremove = &UbuntuMegamSnowflakeRemove{}
	templates.Register("UbuntuMegamSnowflakeRemove", ubuntumegamsnowflakeremove)
}

type UbuntuMegamSnowflakeRemove struct{}

func (tpl *UbuntuMegamSnowflakeRemove) Render(p urknall.Package) {
	p.AddTemplate("snowflake", &UbuntuMegamSnowflakeRemoveTemplate{})
}

func (tpl *UbuntuMegamSnowflakeRemove) Options(opts map[string]string) {
}

func (tpl *UbuntuMegamSnowflakeRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuMegamSnowflakeRemove{})
}

type UbuntuMegamSnowflakeRemoveTemplate struct{}

func (m *UbuntuMegamSnowflakeRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("megamsnowflake",
		Shell("service snowflake stop"),
		RemovePackage("megamsnowflake"),
		RemovePackages(""),
		PurgePackages("megamsnowflake"),
		Shell("dpkg --get-selections megam*"),
	)
}
