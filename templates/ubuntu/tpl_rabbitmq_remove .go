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
	"github.com/dynport/urknall"
	"github.com/megamsys/megdc/templates"
)

var ubunturabbitmqremove *UbuntuRabbitmqRemove

func init() {
	ubunturabbitmqremove = &UbuntuRabbitmqRemove{}
	templates.Register("UbuntuRabbitmqRemove", ubunturabbitmqremove)
}

type UbuntuRabbitmqRemove struct{}

func (tpl *UbuntuRabbitmqRemove) Render(p urknall.Package) {
	p.AddTemplate("rabbitmq", &UbuntuRabbitRemoveTemplate{})
}

func (tpl *UbuntuRabbitmqRemove) Options(opts map[string]string) {
}


func (tpl *UbuntuRabbitmqRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuRabbitmqRemove{})
}

type UbuntuRabbitRemoveTemplate struct{}

func (m *UbuntuRabbitRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("rabbitmq",
		RemovePackage("rabbitmq-server"),
		RemovePackages(""),
		PurgePackages("rabbitmq-server"),
		Shell("dpkg --get-selections rabbitmq*"),
	)
}
