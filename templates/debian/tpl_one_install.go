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
	ONE_INSTALL_LOG = "/var/log/megam/megamcib/opennebula.log"
	Slash = `\`
)

var debianoneinstall *DebianOneInstall

func init() {
	debianoneinstall = &DebianOneInstall{}
	templates.Register("DebianOneInstall", debianoneinstall)
}

type DebianOneInstall struct{}

func (tpl *DebianOneInstall) Render(p urknall.Package) {
	p.AddTemplate("one", &DebianOneInstallTemplate{})
}

func (tpl *DebianOneInstall) Options(t *templates.Template) {}

func (tpl *DebianOneInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianOneInstall{})
}

type DebianOneInstallTemplate struct{}

func (m *DebianOneInstallTemplate) Render(pkg urknall.Package) {

	ip := u.IP("")

	pkg.AddCommands("repository",
		u.Shell("wget -q -O- http://downloads.opennebula.org/repo/Debian/repo.key | apt-key add -"),
		u.Shell("echo 'deb http://downloads.opennebula.org/repo/4.14/Debian/8 stable opennebula' > /etc/apt/sources.list.d/opennebula.list"),
		u.UpdatePackagesOmitError(),
	)

	pkg.AddCommands("one-install",
		u.InstallPackages("build-essential autoconf libtool make lvm2 ssh iproute iputils-arping opennebula opennebula-sunstone ntp ruby-dev"),
	)

	pkg.AddCommands("requires",
		u.Shell("echo 'oneadmin ALL = (root) NOPASSWD:ALL' | tee /etc/sudoers.d/oneadmin"),
		u.Shell(" chmod 0440 /etc/sudoers.d/oneadmin"),
		u.Shell(" rm /usr/share/one/install_gems"),
		u.Shell(" cp /usr/share/megam/megdc/conf/install_gems /usr/share/one/install_gems"),
		u.Shell(" chmod 755 /usr/share/one/install_gems"),
		u.Shell("/usr/share/one/install_gems sunstone"),
		u.Shell("sed -i 's/^[ \t]*:host:.*/:host: "+ip+"/' /etc/one/sunstone-server.conf"),
		u.AsUser("oneadmin",u.Shell("echo 'TM_MAD=ssh' >/tmp/ds_tm_mad")),
		u.AsUser("oneadmin",u.Shell("onedatastore update 0 /tmp/ds_tm_mad")),
		u.AsUser("oneadmin",u.Shell("onedatastore update 1 /tmp/ds_tm_mad")),
		u.AsUser("oneadmin",u.Shell("onedatastore update 2 /tmp/ds_tm_mad")),

		u.Shell("sunstone-server start"),
		u.Shell("econe-server start"),
		u.AsUser("oneadmin",u.Shell("one restart")),
		u.Shell("service opennebula restart"),
	)
}
