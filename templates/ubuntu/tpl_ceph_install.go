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
	CephUser = "megdc"
	User_home = "/home/" + CephUser
	Osd1      = "/storage1"
	Osd2      = "/storage2"
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
osd_pool_default_size = 2
public network = %s/%s
cluster network = %s/%s
mon_pg_warn_max_per_osd = 0
`
)

var ubuntucephinstall *UbuntuCephInstall

func init() {
	ubuntucephinstall = &UbuntuCephInstall{}
	templates.Register("UbuntuCephInstall", ubuntucephinstall)
}

type UbuntuCephInstall struct{}

func (tpl *UbuntuCephInstall) Options(opts map[string]string) {
}

func (tpl *UbuntuCephInstall) Render(p urknall.Package) {
	p.AddTemplate("ceph", &UbuntuCephInstallTemplate{})
}

func (tpl *UbuntuCephInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCephInstall{})
}

type UbuntuCephInstallTemplate struct{}

func (m *UbuntuCephInstallTemplate) Render(pkg urknall.Package) {
	//Host := host()
	Host := ""
	ip := GetLocalIP()

	pkg.AddCommands("sudoer",
		Shell("echo ' "+CephUser+" ALL = (root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/"+CephUser+""),
	)
	pkg.AddCommands("changepermission",
		Shell("sudo chmod 0440 /etc/sudoers.d/"+CephUser+""),
	)

	pkg.AddCommands("ceph_install",
		Shell("sudo echo deb http://ceph.com/debian-hammer/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list"),
		Shell("sudo wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -"),
		Shell("sudo apt-get -y update"),
    InstallPackages("ceph-deploy", "ceph-common", "ceph-mds", "dnsmasq", "openssh-server", "ntp", "sshpass"),
	)

	pkg.AddCommands("edithost",
		Shell("echo '"+ip+" "+Host+"' >> /etc/hosts"),

	)

	pkg.AddCommands("ssh-keygen",
		Mkdir(User_home +"/.ssh",CephUser,0700),
		AsUser(CephUser, Shell("ssh-keygen -N '' -t rsa -f "+User_home+"/.ssh/id_rsa")),
		AsUser(CephUser,	Shell("cp "+User_home+"/.ssh/id_rsa.pub "+User_home+"/.ssh/authorized_keys")),
	)

	pkg.AddCommands("ssh_known_hosts",
		 WriteFile(User_home+"/.ssh/ssh_config", StrictHostKey, CephUser, 0755),
	   WriteFile(User_home+"/.ssh/ssh_config",fmt.Sprintf(SSHHostConfig,Host,Host,CephUser),CephUser, 0755),
	)

	pkg.AddCommands("mkdir_osd",
			Mkdir(Osd1 +"/osd","",0755),
	    Mkdir(Osd2 +"/osd","",0755),
	)

	pkg.AddCommands("ceph-conf",
		AsUser(CephUser, Shell("mkdir "+User_home+"/ceph-cluster")),
			AsUser(CephUser,Shell("cd "+User_home+"/ceph-cluster"),
			AsUser(CephUser,Shell("ceph-deploy new "+Host+" "),
 //add cephconf 	 WriteFile(User_home+"/.ssh/ssh_config", StrictHostKey, CephUser, 0755),
		AsUser(CephUser,Shell("ceph-deploy install "+Host+"")),
		AsUser(CephUser,Shell("ceph-deploy mon create-initial")),
		AsUser(CephUser,Shell("ceph-deploy osd prepare "+Host+":"+Osd1+"/osd "+Host+":"+Osd2+"/osd ")),
		AsUser(CephUser,Shell("ceph-deploy osd activate "+Host+":"+Osd1+"/osd "+Host+":"+Osd2+"/osd ")),
		AsUser(CephUser,Shell("ceph-deploy admin "+Host+"")),
		AsUser(CephUser,Shell("sudo chmod +r /etc/ceph/ceph.client.admin.keyring")),
		AsUser(CephUser,Shell("sleep 180")),
		AsUser(CephUser,Shell("ceph osd pool set rbd pg_num 100")),
		AsUser(CephUser,Shell("sleep 180")),
		AsUser(CephUser,Shell("ceph osd pool set rbd pgp_num 100")),
	)
	pkg.AddCommands("copy",
		Shell("cp "+User_home+"/ceph-cluster/*.keyring /etc/ceph/"),
	)

}
