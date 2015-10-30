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

package templates

import (
	"github.com/dynport/urknall"
)

type megamd struct{}

func (tpl *megamd) Render(p urknall.Package) {
	p.AddTemplate("megamd", &MegamdTemplate{})
}

type MegamdTemplate struct{}

func (m *MegamdTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("build",
		And("ls -la",),
	)
}
