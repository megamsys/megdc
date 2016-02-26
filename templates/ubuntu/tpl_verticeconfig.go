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
	"fmt"
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"
	"github.com/megamsys/megdc/db"
)

var ubuntuverticeconfig *UbuntuVerticeConfig

func init() {
	ubuntuverticeconfig = &UbuntuVerticeConfig{}
	templates.Register("UbuntuVerticeConfig", ubuntuverticeconfig)
}

type UbuntuVerticeConfig struct{
}

func (tpl *UbuntuVerticeConfig) Render(p urknall.Package) {
	p.AddTemplate("vertice-conf", &UbuntuVerticeConfigTemplate{
	})
}

func (tpl *UbuntuVerticeConfig) Options(t *templates.Template) {
}

func (tpl *UbuntuVerticeConfig) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuVerticeConfig{})
}

type UbuntuVerticeConfigTemplate struct{
	path string
}

func (m *UbuntuVerticeConfigTemplate) Render(pkg urknall.Package) {
  c,_ := db.ParseConfig()
  fmt.Println("************************************************")
	fmt.Println(c)
	err := db.StoreDB(c.Common)
	fmt.Println(err)
}
