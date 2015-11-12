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
	ONE_INSTALL_LOG = "/var/log/megam/megamcib/opennebula.log"
	Slash = `\`
)

var ubuntuoneinstall *UbuntuOneInstall

func init() {
	ubuntuoneinstall = &UbuntuOneInstall{}
	templates.Register("UbuntuOneInstall", ubuntuoneinstall)
}

type UbuntuOneInstall struct{}

func (tpl *UbuntuOneInstall) Render(p urknall.Package) {
	p.AddTemplate("one", &UbuntuOneInstallTemplate{})
}

func (tpl *UbuntuOneInstall) Options(opts map[string]string) {
}

func (tpl *UbuntuOneInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuOneInstall{})
}

type UbuntuOneInstallTemplate struct{}

func (m *UbuntuOneInstallTemplate) Render(pkg urknall.Package) {

	ip := IP()

	pkg.AddCommands("repository",
		Shell("wget -q -O- http://downloads.opennebula.org/repo/Ubuntu/repo.key | apt-key add -"),
		Shell("echo 'deb http://downloads.opennebula.org/repo/4.14/Ubuntu/14.04 stable opennebula' > /etc/apt/sources.list.d/opennebula.list"),
		UpdatePackagesOmitError(),
	)

	pkg.AddCommands("one-install",
		InstallPackages("opennebula opennebula-sunstone ntp ruby2.0 ruby2.0-dev ruby-dev"),
	)

	pkg.AddCommands("requires",
		Shell("echo 'oneadmin ALL = (root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/oneadmin"),
		//Shell("sudo apt-get -y install ntp ruby2.0 ruby2.0-dev ruby-dev"),
		Shell("rm /usr/bin/ruby"),
		Shell("rm /usr/bin/gem"),
		Shell("ln -s /usr/bin/ruby2.0 /usr/bin/ruby"),
		Shell("ln -s /usr/bin/gem2.0 /usr/bin/gem"),
		Shell("sudo chmod 0440 /etc/sudoers.d/oneadmin"),
		//Shell("sudo rm /usr/share/one/install_gems"),
		//Shell("sudo cp ~/install_gems /usr/share/one/install_gems"),
		//Shell("sudo cp /usr/share/megam/megdc/conf/trusty/opennebula/install_gems /usr/share/one/install_gems"),
		Shell("sudo chmod 755 /usr/share/one/install_gems"),
		Shell("sudo /usr/share/one/install_gems sunstone"),
		Shell("sed -i 's/^[ \t]*:host:.*/:host: "+ip+"/' /etc/one/sunstone-server.conf"),
		AsUser("oneadmin",Shell("echo 'TM_MAD=ssh' >/tmp/ds_tm_mad")),
		AsUser("oneadmin",Shell("onedatastore update 0 /tmp/ds_tm_mad'")),
		AsUser("oneadmin",Shell("onedatastore update 1 /tmp/ds_tm_mad'")),
		AsUser("oneadmin",Shell("onedatastore update 2 /tmp/ds_tm_mad'")),
		Shell("sunstone-server start"),
		Shell("econe-server start"),

		Shell("sudo -H -u oneadmin bash -c 'one restart'"),
		Shell("service opennebula restart"),
	)
}
