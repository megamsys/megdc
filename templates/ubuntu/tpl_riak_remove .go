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

var ubunturiakremove *UbuntuRiakRemove

func init() {
	ubunturiakremove = &UbuntuRiakRemove{}
	templates.Register("UbuntuRiakRemove", ubunturiakremove)
}

type UbuntuRiakRemove struct{}

func (tpl *UbuntuRiakRemove) Render(p urknall.Package) {
	p.AddTemplate("riak", &UbuntuRiakRemoveTemplate{})
}

func (tpl *UbuntuRiakRemove) Options(opts map[string]string) {
}

func (tpl *UbuntuRiakRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuRiakRemove{})
}

type UbuntuRiakRemoveTemplate struct{}

func (m *UbuntuRiakRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("riak",
		RemovePackage("riak"),
		RemovePackages(""),
		PurgePackages("riak"),
		Shell("dpkg --get-selections riak*"),
	)
}
