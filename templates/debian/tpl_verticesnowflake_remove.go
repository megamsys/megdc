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

var debiansnowflakeremove *DebianSnowflakeRemove

func init() {
	debiansnowflakeremove = &DebianSnowflakeRemove{}
	templates.Register("DebianSnowflakeRemove", debiansnowflakeremove)
}

type DebianSnowflakeRemove struct{}

func (tpl *DebianSnowflakeRemove) Render(p urknall.Package) {
	p.AddTemplate("snowflake", &DebianSnowflakeRemoveTemplate{})
}

func (tpl *DebianSnowflakeRemove) Options(t *templates.Template) {
}

func (tpl *DebianSnowflakeRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianSnowflakeRemove{})
}

type DebianSnowflakeRemoveTemplate struct{}

func (m *DebianSnowflakeRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("verticesnowflake",
		u.Shell("systemctl stop snowflake"),
		u.RemovePackage("verticesnowflake"),
		u.RemovePackages(""),
		u.PurgePackages("verticesnowflake"),
		u.Shell("dpkg --get-selections megam*"),
	)
	pkg.AddCommands("snowflake-clean",
		u.Shell("rm -r /var/lib/urknall/snowflake*"),
	)
}
