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

const (
	KnownHostsList = `
	ConnectTimeout 5
	Host *
	StrictHostKeyChecking no
	`
)

var ubuntusshpass *UbuntuSshPass

func init() {
	ubuntusshpass = &UbuntuSshPass{}
	templates.Register("UbuntuSshPass", ubuntusshpass)
}

type UbuntuSshPass struct{}


func (tpl *UbuntuSshPass) Render(p urknall.Package) {
	p.AddTemplate("onehost", &UbuntuSshPassTemplate{})
}

func (tpl *UbuntuSshPass) Options(opts map[string]string) {

}


func (tpl *UbuntuSshPass) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuSshPass{})
}

type UbuntuSshPassTemplate struct{}

func (m *UbuntuSshPassTemplate) Render(pkg urknall.Package) {

	ip := IP()

	pkg.AddCommands("install-sshpass",
	  InstallPackages("sshpass"),
	)
	pkg.AddCommands("SSHPass",
		Shell("sudo -H -u oneadmin bash -c 'sshpass -p "+Slash+"'oneadmin'"+Slash+" scp -o StrictHostKeyChecking=no /var/lib/one/.ssh/id_rsa.pub oneadmin@"+ ip +":/var/lib/one/.ssh/authorized_keys'"),
    WriteFile("/var/lib/one/.ssh/config",KnownHostsList,"oneadmin", 0755),
	)
}
