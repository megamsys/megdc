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
	"github.com/megamsys/urknall"
	"github.com/megamsys/megdc/templates"
	u "github.com/megamsys/megdc/templates/ubuntu"
)

const (
	KnownHostsList = `
	ConnectTimeout 5
	Host *
	StrictHostKeyChecking no
	`
)

var debiansshpass *DebianSshPass

func init() {
	debiansshpass = &DebianSshPass{}
	templates.Register("DebianSshPass", debiansshpass)
}

type DebianSshPass struct{
	Host string
}


func (tpl *DebianSshPass) Render(p urknall.Package) {
	p.AddTemplate("sshpass", &DebianSshPassTemplate{
		Host: tpl.Host,
	})
}

func (tpl *DebianSshPass) Options(t *templates.Template) {

if hs, ok := t.Options["HOST"]; ok {
	tpl.Host = hs
}
}


func (tpl *DebianSshPass) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianSshPass{
		Host: tpl.Host,
	})
}

type DebianSshPassTemplate struct{
	Host string
}

func (m *DebianSshPassTemplate) Render(pkg urknall.Package) {

	pkg.AddCommands("install-sshpass",
	  u.InstallPackages("sshpass"),
	)
	pkg.AddCommands("SSHPass",
		u.AsUser("oneadmin", u.Shell("sshpass -p 'oneadmin' scp -o StrictHostKeyChecking=no /var/lib/one/.ssh/id_rsa.pub oneadmin@"+ m.Host +":/var/lib/one/.ssh/authorized_keys")),
    u.WriteFile("/var/lib/one/.ssh/config",KnownHostsList,"oneadmin", 0755),
	)
}
