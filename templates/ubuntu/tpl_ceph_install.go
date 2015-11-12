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
	"github.com/megamsys/urknall"
)

const (
	CephUser = "CephUser"
	Osd1     = "Osd1"
	Osd2     = "Osd2"

	UserHomePrefix = "/home/"

	StrictHostKey = `#!/bin/sh

	ConnectTimeout 5
	Host *
	StrictHostKeyChecking no
	`

	SSHHostConfig = `#!/bin/sh
  Host %s
 Hostname %s
 User %s
`
	CephConf = `osd crush chooseleaf type = 0
osd_pool_default_size = %s
public network = %s
cluster network = %s
mon_pg_warn_max_per_osd = 0
`
)

var ubuntucephinstall *UbuntuCephInstall

func init() {
	ubuntucephinstall = &UbuntuCephInstall{}
	templates.Register("UbuntuCephInstall", ubuntucephinstall)
}

type UbuntuCephInstall struct {
	osd1     string
	osd2     string
	cephuser string
}

func (tpl *UbuntuCephInstall) Options(opts map[string]string) {

		if osd1, ok := opts[Osd1]; ok {
		tpl.osd1 = osd1
	}
	if osd2, ok := opts[Osd2]; ok {
		tpl.osd2 = osd2
	}
	if cephuser, ok := opts[CephUser]; ok {
		tpl.cephuser = cephuser
	}
}

func (tpl *UbuntuCephInstall) Render(p urknall.Package) {
	p.AddTemplate("ceph", &UbuntuCephInstallTemplate{
		osd1:     tpl.osd1,
		osd2:     tpl.osd2,
		cephuser: tpl.cephuser,
		cephhome: UserHomePrefix + tpl.cephuser,
	})
}

func (tpl *UbuntuCephInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCephInstall{
		osd1:     tpl.osd1,
		osd2:     tpl.osd2,
		cephuser: tpl.cephuser,

	})
}

type UbuntuCephInstallTemplate struct {
	osd1     string
	osd2     string
	cephuser string
	cephhome string
}

func (m *UbuntuCephInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := IP()
	Osd1 := m.osd1
	Osd2 := m.osd2
	CephUser := m.cephuser
	CephHome := m.cephhome
	pkg.AddCommands("cephuser_sudoer",
		Shell("echo '"+CephUser+" ALL = (root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/"+CephUser+""),
	)
	pkg.AddCommands("chmod_sudoer",
		Shell("sudo chmod 0440 /etc/sudoers.d/"+CephUser+""),
	)

	pkg.AddCommands("cephinstall",
		Shell("sudo echo deb http://ceph.com/debian-hammer/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list"),
		Shell("sudo wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -"),
		Shell("sudo apt-get -y update"),
		InstallPackages("ceph-deploy", "ceph-common", "ceph-mds", "dnsmasq", "openssh-server", "ntp", "sshpass"),
	)

	pkg.AddCommands("getip",
		Shell("ip3=`echo 103.56.92.24| cut -d'.' -f 1,2,3`"),
)
	pkg.AddCommands("etchost",
		Shell("echo '"+ip+" "+host+"' >> /etc/hosts"),
	)

	pkg.AddCommands("ssh-keygen",
		Mkdir(CephHome+"/.ssh", CephUser, 0700),
		AsUser(CephUser, Shell("ssh-keygen -N '' -t rsa -f "+CephHome+"/.ssh/id_rsa")),
		AsUser(CephUser, Shell("cp "+CephHome+"/.ssh/id_rsa.pub "+CephHome+"/.ssh/authorized_keys")),

	)

	pkg.AddCommands("ssh_known_hosts",
		WriteFile(CephHome+"/.ssh/ssh_config", StrictHostKey, CephUser, 0755),
		WriteFile(CephHome+"/.ssh/ssh_config", fmt.Sprintf(SSHHostConfig, host, host, CephUser), CephUser, 0755),
	)

	pkg.AddCommands("mkdir_osd",
		Mkdir(Osd1+"/osd", "", 0755),
		Mkdir(Osd2+"/osd", "", 0755),
	)

	pkg.AddCommands("write_cephconf",
		AsUser(CephUser, Shell("mkdir "+CephHome+"/ceph-cluster")),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster")),
		AsUser(CephUser, Shell("ceph-deploy new "+host+" ")),
		WriteFile(CephHome+"/ceph-cluster/ceph.conf",
			fmt.Sprintf(CephConf, m.osdPoolSize(Osd1, Osd2), m.slashIp(), m.slashIp()), CephUser, 0755),

		AsUser(CephUser, Shell("ceph-deploy install "+host+"")),
		AsUser(CephUser, Shell("ceph-deploy mon create-initial")),
		AsUser(CephUser, Shell("ceph-deploy osd prepare "+host+":"+Osd1+"/osd "+host+":"+Osd2+"/osd ")),
		AsUser(CephUser, Shell("ceph-deploy osd activate "+host+":"+Osd1+"/osd "+host+":"+Osd2+"/osd ")),
		AsUser(CephUser, Shell("ceph-deploy admin "+host+"")),
		AsUser(CephUser, Shell("sudo chmod +r /etc/ceph/ceph.client.admin.keyring")),
		AsUser(CephUser, Shell("sleep 180")),
		AsUser(CephUser, Shell("ceph osd pool set rbd pg_num 100")),
		AsUser(CephUser, Shell("sleep 180")),
		AsUser(CephUser, Shell("ceph osd pool set rbd pgp_num 100")),
	)
	pkg.AddCommands("copy_keyring",
		Shell("cp "+CephHome+"/ceph-cluster/*.keyring /etc/ceph/"),
	)
}

func (m *UbuntuCephInstallTemplate) noOfIpsFromMask() int {
	si, _ := IPNet().Mask.Size() //from your netwwork
	return si
}

func (m *UbuntuCephInstallTemplate) slashIp() string {
	s := strings.Split(IP(), ".")
	p := s[0 : len(s)-1]
	p = append(p, "0")
	return fmt.Sprintf("%s/%d", strings.Join(p, "."), m.noOfIpsFromMask())
}

func (m *UbuntuCephInstallTemplate) osdPoolSize(osds ...string) int {
	return len(osds)
}
