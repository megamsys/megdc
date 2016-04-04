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
	"github.com/megamsys/urknall"
	"github.com/megamsys/megdc/templates"

)

const (
	KnownHostsList = `
	ConnectTimeout 5
	Host *
	StrictHostKeyChecking no
	`
)

var centossshpass *CentOsSshPass

func init() {
	centossshpass = &CentOsSshPass{}
	templates.Register("CentosSshPass", centossshpass)
}

type CentOsSshPass struct{
	Host string
}


func (tpl *CentOsSshPass) Render(p urknall.Package) {
	p.AddTemplate("sshpass", &CentOsSshPassTemplate{
		Host: tpl.Host,
	})
}

func (tpl *CentOsSshPass) Options(t *templates.Template) {

if hs, ok := t.Options["HOST"]; ok {
	tpl.Host = hs
}
}


func (tpl *CentOsSshPass) Run(target urknall.Target) error {
	return urknall.Run(target, &CentOsSshPass{
		Host: tpl.Host,
	})
}

type CentOsSshPassTemplate struct{
	Host string
}

func (m *CentOsSshPassTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("install-sshpass",
	  InstallPackages("sshpass"),
	)
	pkg.AddCommands("SSHPass",
		AsUser("oneadmin", Shell("sshpass -p 'oneadmin' scp -o StrictHostKeyChecking=no /var/lib/one/.ssh/id_rsa.pub oneadmin@"+ m.Host +":/var/lib/one/.ssh/authorized_keys")),
    WriteFile("/var/lib/one/.ssh/config",KnownHostsList,"oneadmin", 0755),
	)
}
