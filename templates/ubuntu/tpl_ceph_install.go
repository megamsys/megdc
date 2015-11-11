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
	User_home = "/home/megdc"
	Osd1      = "/storage1"
	Osd2      = "/storage2"
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
	pkg.AddCommands("cephuser",
		Shell(" echo 'Make ceph user as sudoer'"),
	)
	pkg.AddCommands("sudoer",
		Shell("echo ' "+CephUser+" ALL = (root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/"+CephUser+""),
	)
	pkg.AddCommands("changepermission",
		Shell("sudo chmod 0440 /etc/sudoers.d/"+CephUser+""),
	)
	pkg.AddCommands("startinstall",
		Shell("echo 'Started installing ceph'"),
	)
	pkg.AddCommands("install",
		Shell("sudo echo deb http://ceph.com/debian-hammer/ $(lsb_release -sc) main | sudo tee /etc/apt/sources.list.d/ceph.list"),
	)
	pkg.AddCommands("get",
		Shell("sudo wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | sudo apt-key add -"),
	)
	pkg.AddCommands("update",
		Shell("sudo apt-get -y update"),
	)
	pkg.AddCommands("cephDeployinstall",

		InstallPackages("ceph-deploy", "ceph-common", "ceph-mds", "dnsmasq", "openssh-server", "ntp", "sshpass"),
	)

	pkg.AddCommands("ipaddress",
		Shell("IP_ADDR="+ip+""),
	)
	pkg.AddCommands("entry",
		Shell("echo 'Adding entry in /etc/hosts'"),
	)
	pkg.AddCommands("edithost",
		Shell("echo '"+ip+" "+Host+"'"),
	)
	pkg.AddCommands("ssh",
		Shell("echo 'Processing ssh-keygen'"),
	)
	pkg.AddCommands("adduser",
		AddUser("megdc", true),
		Shell("ssh-keygen -N '' -t rsa -f "+User_home+"/.ssh/id_rsa"),
		Shell("cp "+User_home+"/.ssh/id_rsa.pub "+User_home+"/.ssh/authorized_keys"),
	)
	pkg.AddCommands("ipKnown_hosts",
		AddUser("megdc", true),
		WriteFile(""+User_home+"/.ssh/ssh_config", content, ""+CephUser+"", 0755),
	)
	pkg.AddCommands("hostuser",
		AddUser(""+CephUser+"", true),
		WriteFile(""+User_home+"/.ssh/config", content2, ""+CephUser+"", 0755),
	)
	pkg.AddCommands("makeosd",
		Shell("echo 'Making directory inside osd drive '"),
	)

	pkg.AddCommands("osd1",
		Shell("mkdir "+Osd1+"/osd"),
	)
	pkg.AddCommands("osd2",
		Shell("mkdir "+Osd2+"/osd"),
	)
	pkg.AddCommands("getip",
		Shell("ip3=`echo 103.56.92.24| cut -d'.' -f 1,2,3`"),
	)

	pkg.AddCommands("cephconfig",
		Shell("echo 'Ceph configuration started...'"),
	)
	pkg.AddCommands("conf",
		AddUser(""+CephUser+"", true),
		Shell("mkdir "+User_home+"/ceph-cluster"),
		Shell("cd "+User_home+"/ceph-cluster"),
		Shell("ceph-deploy new "+Host+" "),
		Shell("echo 'osd crush chooseleaf type = 0'"),
		Shell("echo 'public network = $ip3.0/24'"),
		Shell("echo 'cluster network = $ip3.0/24'"),
		Shell("ceph-deploy install "+Host+""),
		Shell("ceph-deploy mon create-initial"),
		Shell("ceph-deploy osd prepare "+Host+":"+Osd1+"/osd "+Host+":"+Osd2+"/osd "),
		Shell("ceph-deploy osd activate "+Host+":"+Osd1+"/osd "+Host+":"+Osd2+"/osd "),
		Shell("ceph-deploy admin "+Host+""),
		Shell("sudo chmod +r /etc/ceph/ceph.client.admin.keyring"),
		Shell("sleep 180"),
		Shell("ceph osd pool set rbd pg_num 150"),
		Shell("sleep 180"),
		Shell("ceph osd pool set rbd pgp_num 150"),
	)
	pkg.AddCommands("copy",
		Shell("cp "+User_home+"/ceph-cluster/*.keyring /etc/ceph/"),
	)
	pkg.AddCommands("complete",
		Shell("echo 'Ceph installed successfully.'"),
	)

}

const content = `#!/bin/sh

ConnectTimeout 5
Host *
StrictHostKeyChecking no
`
const content2 = `#!/bin/sh
  Host ranjitha-sfd-sdf
 Hostname ranjitha-sfd-sdf
 User megdc

`
