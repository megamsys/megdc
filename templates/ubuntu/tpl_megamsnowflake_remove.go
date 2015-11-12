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

var ubuntusnowflakeremove *UbuntuSnowflakeRemove

func init() {
	ubuntusnowflakeremove = &UbuntuSnowflakeRemove{}
	templates.Register("UbuntuSnowflakeRemove", ubuntusnowflakeremove)
}

type UbuntuSnowflakeRemove struct{}

func (tpl *UbuntuSnowflakeRemove) Render(p urknall.Package) {
	p.AddTemplate("snowflake", &UbuntuSnowflakeRemoveTemplate{})
}

func (tpl *UbuntuSnowflakeRemove) Options(opts map[string]string) {
}

func (tpl *UbuntuSnowflakeRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuSnowflakeRemove{})
}

type UbuntuSnowflakeRemoveTemplate struct{}

func (m *UbuntuSnowflakeRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("megamsnowflake",
		Shell("service snowflake stop"),
		RemovePackage("megamsnowflake"),
		RemovePackages(""),
		PurgePackages("megamsnowflake"),
		Shell("dpkg --get-selections megam*"),
	)
}
