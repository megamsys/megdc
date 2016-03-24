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
	"fmt"
	"os"
	"strings"
	"github.com/megamsys/megdc/templates"
	u "github.com/megamsys/megdc/templates/ubuntu"
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

var debiancephinstall *DebianCephInstall

func init() {
	debiancephinstall = &DebianCephInstall{}
	templates.Register("DebianCephInstall", debiancephinstall)
}

type DebianCephInstall struct {
	osds      []string
	cephuser string
	phydev    string
}

func (tpl *DebianCephInstall) Options(t *templates.Template) {
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

func (tpl *DebianCephInstall) Render(p urknall.Package) {
	p.AddTemplate("ceph", &DebianCephInstallTemplate{
		osds:     tpl.osds,
		cephuser: tpl.cephuser,
		cephhome: UserHomePrefix + tpl.cephuser,
		phydev:    tpl.phydev,
	})
}

func (tpl *DebianCephInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianCephInstall{
		osds:     tpl.osds,
		cephuser: tpl.cephuser,
		phydev:    tpl.phydev,

	})
}

type DebianCephInstallTemplate struct {
  osds     []string
	cephuser string
	cephhome string
	phydev    string
}

func (m *DebianCephInstallTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	ip := u.IP(m.phydev)
  osddir := u.ArraytoString("/","/osd",m.osds)
	hostosd := u.ArraytoString(host+":/","/osd",m.osds)
	CephUser := m.cephuser
	CephHome := m.cephhome
	
 pkg.AddCommands("cephinstall",
		u.Shell("echo deb https://download.ceph.com/debian-infernalis/ jessie main | tee /etc/apt/sources.list.d/ceph.list"),
		u.Shell("wget -q -O- 'https://download.ceph.com/keys/release.asc' | apt-key add -"),
		u.InstallPackages("apt-transport-https  sudo"),
	  u.UpdatePackagesOmitError(),
		u.InstallPackages("ceph-deploy ceph-common ceph-mds dnsmasq openssh-server ntp sshpass ceph ceph-mds ceph-deploy radosgw"),
	)

	pkg.AddCommands("cephuser_add",
	 u.AddUser(CephUser,false),
	)

	pkg.AddCommands("cephuser_sudoer",
		u.Shell("echo '"+CephUser+" ALL = (root) NOPASSWD:ALL' && mkdir -p /etc/sudoers.d  | tee /etc/sudoers.d/"+CephUser+""),
	)

	pkg.AddCommands("chmod_sudoer",
		u.Shell("chmod 0440 /etc/sudoers.d/"+CephUser+""),
	)

	pkg.AddCommands("etchost",
		u.Shell("echo '"+ip+" "+host+"' >> /etc/hosts"),
	)

	pkg.AddCommands("ssh-keygen",
		u.Mkdir(CephHome+"/.ssh", CephUser, 0700),
		u.AsUser(CephUser, u.Shell("ssh-keygen -N '' -t rsa -f "+CephHome+"/.ssh/id_rsa")),
		u.AsUser(CephUser, u.Shell("cp "+CephHome+"/.ssh/id_rsa.pub "+CephHome+"/.ssh/authorized_keys")),
	)

	pkg.AddCommands("ssh_known_hosts",
		u.WriteFile(CephHome+"/.ssh/ssh_config", StrictHostKey, CephUser, 0755),
		u.WriteFile(CephHome+"/.ssh/config", fmt.Sprintf(SSHHostConfig, host, host, CephUser), CephUser, 0755),
	)

	pkg.AddCommands("mkdir_osd",
		u.Mkdir(osddir,"", 0755),
		u.Shell("sudo chown -R"+CephUser+":"+CephUser+" "+osddir ),
	)

	pkg.AddCommands("write_cephconf",
		u.AsUser(CephUser, u.Shell("mkdir "+CephHome+"/ceph-cluster")),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster")),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy new "+host+" ")),
	  u.AsUser(CephUser, u.Shell("echo 'osd crush chooseleaf type = 0' >> "+CephHome+"/ceph-cluster/ceph.conf")),
		u.AsUser(CephUser,u.Shell("echo 'osd_pool_default_size = 2' >> "+CephHome+"/ceph-cluster/ceph.conf")),
		u.AsUser(CephUser,u.Shell("echo 'mon_pg_warn_max_per_osd = 0' >> "+CephHome+"/ceph-cluster/ceph.conf")),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy install "+host+"")),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy mon create-initial | ceph-deploy mon create-initial")),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy osd prepare "+ hostosd )),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy osd activate "+ hostosd )),
		u.AsUser(CephUser, u.Shell("cd "+CephHome+"/ceph-cluster;ceph-deploy admin "+host+"")),
		u.AsUser(CephUser, u.Shell("chmod +r /etc/ceph/ceph.client.admin.keyring")),
		u.AsUser(CephUser, u.Shell("sleep 180")),
		u.AsUser(CephUser, u.Shell("ceph osd pool set rbd pg_num 150")),
		u.AsUser(CephUser, u.Shell("sleep 180")),
		u.AsUser(CephUser, u.Shell("ceph osd pool set rbd pgp_num 150")),
	)
	pkg.AddCommands("copy_keyring",
		u.Shell("cp "+CephHome+"/ceph-cluster/*.keyring /etc/ceph/"),
	)
}

func (m *DebianCephInstallTemplate) noOfIpsFromMask() int {
	si, _ := u.IPNet(m.phydev).Mask.Size() //from your netwwork
	return si
}

func (m *DebianCephInstallTemplate) slashIp() string {
	s := strings.Split(u.IP(m.phydev), ".")
	p := s[0 : len(s)-1]
	p = append(p, "0")
	return fmt.Sprintf("%s/%d", strings.Join(p, "."), m.noOfIpsFromMask())
}

func (m *DebianCephInstallTemplate) osdPoolSize(osds ...string) int {
	return len(osds)
}
