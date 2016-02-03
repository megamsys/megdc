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
	ONEHOST_INSTALL_LOG = "/var/log/megam/megamcib/opennebulahost.log"
)

var debianonehostinstall *DebianOneHostInstall

func init() {
	debianonehostinstall = &DebianOneHostInstall{}
	templates.Register("DebianOneHostInstall", debianonehostinstall)
}

type DebianOneHostInstall struct{}

func (tpl *DebianOneHostInstall) Render(p urknall.Package) {
	p.AddTemplate("onehost", &DebianOneHostInstallTemplate{})
}

func (tpl *DebianOneHostInstall) Options(t *templates.Template) {
}

func (tpl *DebianOneHostInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianOneHostInstall{})
}

type DebianOneHostInstallTemplate struct{}

func (m *DebianOneHostInstallTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("repository",
		u.InstallPackages("sudo"),
		u.Shell("wget -q -O- http://downloads.opennebula.org/repo/Debian/repo.key | apt-key add -"),
		u.Shell("echo 'deb http://downloads.opennebula.org/repo/4.14/Debian/8 stable opennebula' > /etc/apt/sources.list.d/opennebula.list"),
		u.UpdatePackagesOmitError(),
	)
	pkg.AddCommands("depends",
		u.InstallPackages("qemu-system-x86 qemu-kvm libvirt-bin build-essential genromfs autoconf libtool qemu-utils libvirt0 bridge-utils lvm2 ssh iproute iputils-arping make"),
	)

	pkg.AddCommands("one-node",
		u.InstallPackages("opennebula-node"),
	)
  pkg.AddCommands("node",
		u.Shell("sudo usermod -p $(echo oneadmin | openssl passwd -1 -stdin) oneadmin"),
	)
	pkg.AddCommands("vswitch",
		u.InstallPackages("openvswitch-common openvswitch-switch bridge-utils"),
	)

}
