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
	//"github.com/megamsys/libgo/cmd"
)

const (
	CephUser = "CephUser"
	Osd     = "Osd"
	Phydev    = "PhyDev"
	UserHomePrefix = "/home/"

	StrictHostKey = `
	ConnectTimeout 5
	Host *
	StrictHostKeyChecking no
	`

	SSHHostConfig = `
Host %s
 Hostname %s
 User %s
`
	CephConf = `osd crush chooseleaf type = 0
osd_pool_default_size = %d
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
	osds      []string
	cephuser string
	phydev    string
}

func (tpl *UbuntuCephInstall) Options(t *templates.Template) {
	if osds, ok := t.Maps[Osd]; ok {
		tpl.osds = osds
	}
	if cephuser, ok := t.Options[CephUser]; ok {
		tpl.cephuser = cephuser
	}
	if phydev, ok := t.Options[Phydev]; ok {
		tpl.phydev = phydev
	}
}

func (tpl *UbuntuCephInstall) Render(p urknall.Package) {
	p.AddTemplate("ceph", &UbuntuCephInstallTemplate{
		osds:     tpl.osds,
		cephuser: tpl.cephuser,
		cephhome: UserHomePrefix + tpl.cephuser,
		phydev:    tpl.phydev,
	})
}

func (tpl *UbuntuCephInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCephInstall{
		osds:     tpl.osds,
		cephuser: tpl.cephuser,
		phydev:    tpl.phydev,

	})
}

type UbuntuCephInstallTemplate struct {
  osds     []string
	cephuser string
	cephhome string
	phydev    string
}

func (m *UbuntuCephInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := IP(m.phydev)
  osddir := ArraytoString("/","/osd",m.osds)
	hostosd := ArraytoString(host+":/","/osd",m.osds)
	CephUser := m.cephuser
	CephHome := m.cephhome

	pkg.AddCommands("cephinstall",
		 Shell("echo deb https://download.ceph.com/debian-infernalis/ $(lsb_release -sc) main | tee /etc/apt/sources.list.d/ceph.list"),
		 Shell("wget -q -O- 'https://download.ceph.com/keys/release.asc' | apt-key add -"),
		 InstallPackages("apt-transport-https  sudo"),
		 UpdatePackagesOmitError(),
		 InstallPackages("ceph-deploy ceph-common ceph-mds dnsmasq openssh-server ntp sshpass ceph ceph-mds ceph-deploy radosgw"),
	 )

	 pkg.AddCommands("cephuser_add",
 	 AddUser(CephUser,false),
 	)
 	pkg.AddCommands("cephuser_sudoer",
 		Shell("echo '"+CephUser+" ALL = (root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/"+CephUser+""),
 	)
 	pkg.AddCommands("chmod_sudoer",
 		Shell("sudo chmod 0440 /etc/sudoers.d/"+CephUser+""),
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
		WriteFile(CephHome+"/.ssh/config", fmt.Sprintf(SSHHostConfig, host, host, CephUser), CephUser, 0755),
	)

	pkg.AddCommands("mkdir_osd",
		Mkdir(osddir,"", 0755),
		Shell("sudo chown -R "+CephUser+":"+CephUser+" "+osddir ),
	)

	pkg.AddCommands("write_cephconf",
		AsUser(CephUser, Shell("mkdir "+CephHome+"/ceph-cluster")),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster")),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy new "+host+" ")),
	  	AsUser(CephUser, Shell("echo 'osd crush chooseleaf type = 0' >> "+CephHome+"/ceph-cluster/ceph.conf")),
			AsUser(CephUser,Shell("echo 'osd_pool_default_size = 2' >> "+CephHome+"/ceph-cluster/ceph.conf")),
		AsUser(CephUser,Shell("echo 'mon_pg_warn_max_per_osd = 0' >> "+CephHome+"/ceph-cluster/ceph.conf")),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy install "+host+"")),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy mon create-initial")),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy osd prepare "+ hostosd )),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy osd activate "+ hostosd )),
		AsUser(CephUser, Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy admin "+host+"")),
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
	si, _ := IPNet(m.phydev).Mask.Size() //from your netwwork
	return si
}

func (m *UbuntuCephInstallTemplate) slashIp() string {
	s := strings.Split(IP(m.phydev), ".")
	p := s[0 : len(s)-1]
	p = append(p, "0")
	return fmt.Sprintf("%s/%d", strings.Join(p, "."), m.noOfIpsFromMask())
}

func (m *UbuntuCephInstallTemplate) osdPoolSize(osds ...string) int {
	return len(osds)
}
