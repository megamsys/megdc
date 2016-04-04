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

/*
import (
	"os"
	"fmt"
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"
	//"github.com/megamsys/libgo/cmd"
)

const (
	Bridge = "Bridge"
	Hdd     = "Osd"
	Phy    = "PhyDev"
	VgName  = "VgName"
)

var ubuntulvminstall *UbuntuLvmInstall

func init() {
	ubuntulvminstall = &UbuntuLvmInstall{}
	templates.Register("UbuntuLvmInstall", ubuntulvminstall)
}

type UbuntuLvmInstall struct {
	osds      []string
	bridge string
	phydev    string
	vgname string
}

func (tpl *UbuntuLvmInstall) Options(t *templates.Template) {
	if osds, ok := t.Maps[Hdd]; ok {
		tpl.osds = osds
	}
	if bridge, ok := t.Options[Bridge]; ok {
		tpl.bridge = bridge
	}
	if phydev, ok := t.Options[Phy]; ok {
		tpl.phydev = phydev
	}
	if vgname, ok := t.Options[VgName]; ok {
		tpl.vgname = vgname
	}
}

func (tpl *UbuntuLvmInstall) Render(p urknall.Package) {
	p.AddTemplate("lvm", &UbuntuLvmInstallTemplate{
		osds:     tpl.osds,
		bridge: tpl.bridge,
	  vgname: tpl.vgname,
		phydev:    tpl.phydev,
	})
}

func (tpl *UbuntuLvmInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuLvmInstall{
		osds:     tpl.osds,
		bridge: tpl.bridge,
		phydev:    tpl.phydev,
		vgname: tpl.vgname,

	})
}

type UbuntuLvmInstallTemplate struct {
  osds     []string
	bridge string
	vgname string
	phydev    string
}

func (m *UbuntuLvmInstallTemplate) Render(pkg urknall.Package) {
	host,_ := os.Hostname()
	phy := m.phydev
	ip := IP(phy)
  osddir := ArraytoString("/dev/","",m.osds)
	bridge := m.bridge
	vg := m.vgname
  fmt.Println("bridge  ",bridge )
 pkg.AddCommands("lvminstall",
	  UpdatePackagesOmitError(),
		InstallPackages("clvm lvm2 kvm libvirt-bin ruby nfs-common bridge-utils"),
	)
	pkg.AddCommands("vg-setup",
		Shell("ip addr flush dev "+phy+""),
		Shell("brctl addbr "+ bridge),
		Shell("brctl addif "+ bridge+" "+phy+""),
		Shell("pvcreate "+osddir+""),
		Shell("vgcreate "+vg+" "+osddir+""),
	)
}
*/
