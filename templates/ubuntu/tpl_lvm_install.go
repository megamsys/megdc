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

var ubuntulvminstall *UbuntuLvmInstall

func init() {
	ubuntulvminstall = &UbuntuLvmInstall{}
	templates.Register("UbuntuLvmInstall", ubuntulvminstall)
}

type UbuntuLvmInstall struct {
	osds      []string
	bridge string
	phydev    string
}

func (tpl *UbuntuLvmInstall) Options(t *templates.Template) {
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

func (tpl *UbuntuLvmInstall) Render(p urknall.Package) {
	p.AddTemplate("lvm", &UbuntuLvmInstallTemplate{
		osds:     tpl.osds,
		bridge: tpl.bridge,
		lvmhome: UserHomePrefix + tpl.bridge,
		phydev:    tpl.phydev,
	})
}

func (tpl *UbuntuLvmInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuLvmInstall{
		osds:     tpl.osds,
		bridge: tpl.bridge,
		phydev:    tpl.phydev,

	})
}

type UbuntuLvmInstallTemplate struct {
  osds     []string
	bridge string
	lvmhome string
	phydev    string
}

func (m *UbuntuLvmInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := u.IP(m.phydev)
  osddir := u.ArraytoString("/","/osd",m.osds)
	hostosd := u.ArraytoString(host+":/","/osd",m.osds)
	LvmUser := m.bridge
	LvmHome := m.lvmhome

 pkg.AddCommands("lvminstall",
	  u.UpdatePackagesOmitError(),
		u.InstallPackages("lvm lvm2"),
	)

	pkg.AddCommands("bridge_add",
	 u.AddUser(LvmUser,false),
	)


	)
}

func (m *UbuntuLvmInstallTemplate) noOfIpsFromMask() int {
	si, _ := u.IPNet(m.phydev).Mask.Size() //from your netwwork
	return si
}

func (m *UbuntuLvmInstallTemplate) slashIp() string {
	s := strings.Split(u.IP(m.phydev), ".")
	p := s[0 : len(s)-1]
	p = append(p, "0")
	return fmt.Sprintf("%s/%d", strings.Join(p, "."), m.noOfIpsFromMask())
}

func (m *UbuntuLvmInstallTemplate) osdPoolSize(osds ...string) int {
	return len(osds)
}
