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
/*
import (
	"fmt"
	"os"
	"strings"
	"github.com/megamsys/megdc/templates"
	u "github.com/megamsys/megdc/templates/ubuntu"
	"github.com/megamsys/urknall"
	//"github.com/megamsys/libgo/cmd"
)

const (
	Bridge = "Bridge"
	Osd     = "Osd"
	Phydev    = "PhyDev"
	UserHomePrefix = "/home/"
)

var debianlvminstall *DebianLvmInstall

func init() {
	debianlvminstall = &DebianLvmInstall{}
	templates.Register("DebianLvmInstall", debianlvminstall)
}

type DebianLvmInstall struct {
	osds      []string
	bridge string
	phydev    string
}

func (tpl *DebianLvmInstall) Options(t *templates.Template) {
	if osds, ok := t.Maps[Osd]; ok {
		tpl.osds = osds
	}
	if bridge, ok := t.Options[LvmUser]; ok {
		tpl.bridge = bridge
	}
	if phydev, ok := t.Options[Phydev]; ok {
		tpl.phydev = phydev
	}
}

func (tpl *DebianLvmInstall) Render(p urknall.Package) {
	p.AddTemplate("lvm", &DebianLvmInstallTemplate{
		osds:     tpl.osds,
		bridge: tpl.bridge,
		lvmhome: UserHomePrefix + tpl.bridge,
		phydev:    tpl.phydev,
	})
}

func (tpl *DebianLvmInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianLvmInstall{
		osds:     tpl.osds,
		bridge: tpl.bridge,
		phydev:    tpl.phydev,

	})
}

type DebianLvmInstallTemplate struct {
  osds     []string
	bridge string
	lvmhome string
	phydev    string
}

func (m *DebianLvmInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := u.IP(m.phydev)
  osddir := u.ArraytoString("/dev","/",m.osds)
	LvmBride := m.bridge
	LvmHome := m.lvmhome

 pkg.AddCommands("lvminstall",
	  u.UpdatePackagesOmitError(),
		u.InstallPackages("clvm lvm2 kvm libvirt-bin ruby nfs-common bridge-utils"),
	)
	pkg.AddCommands("lvminstall",
 	  u.Shell("echo "+ osddir +">test")
 	)

}

func (m *DebianLvmInstallTemplate) noOfIpsFromMask() int {
	si, _ := u.IPNet(m.phydev).Mask.Size() //from your netwwork
	return si
}

func (m *DebianLvmInstallTemplate) slashIp() string {
	s := strings.Split(u.IP(m.phydev), ".")
	p := s[0 : len(s)-1]
	p = append(p, "0")
	return fmt.Sprintf("%s/%d", strings.Join(p, "."), m.noOfIpsFromMask())
}

func (m *DebianLvmInstallTemplate) osdPoolSize(osds ...string) int {
	return len(osds)
}
*/
