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

package centos

import (
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"
)

var centosverticecommonremove *CentOsMegamCommonRemove

func init() {
	centosverticecommonremove = &CentOsMegamCommonRemove{}
	templates.Register("CentOsMegamCommonRemove", centosverticecommonremove)
}

type CentOsMegamCommonRemove struct{}

func (tpl *CentOsMegamCommonRemove) Render(p urknall.Package) {
	p.AddTemplate("common", &CentOsMegamCommonRemoveTemplate{})
}

func (tpl *CentOsMegamCommonRemove) Options(t *templates.Template) {
}

func (tpl *CentOsMegamCommonRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &CentOsMegamCommonRemove{})
}

type CentOsMegamCommonRemoveTemplate struct{}

func (m *CentOsMegamCommonRemoveTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("verticecommon",
		RemovePackage("verticecommon"),
		RemovePackages(""),
		Shell("dpkg --get-selections megam*"),
	)
	pkg.AddCommands("common-clean",
		Shell("rm -rf /var/lib/urknall/common*"),
	)
}
