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
	"github.com/dynport/urknall"
	"github.com/megamsys/megdc/templates"
)

const (
	// DefaultCephRepo is the default megam repository to install if its not provided.
	Ceph_user = "cibadmin"
	Host      = `hostname`
	User_home = "/home/cibadmin"
	Osd1      = "/storage1"
	Osd2      = "/storage2"
	Osd3      = "/storage3"
)

var ubuntucephinstall *UbuntuCephInstall

func init() {
	ubuntucephinstall = &UbuntuCephInstall{}
	templates.Register("UbuntuCephInstall", ubuntucephinstall)
}

type UbuntuCephInstall struct{}

func (tpl *UbuntuCephInstall) Render(p urknall.Package) {
	p.AddTemplate("ceph", &UbuntuCephInstallTemplate{})
}

func (tpl *UbuntuCephInstall) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCephInstall{})
}

type UbuntuCephInstallTemplate struct{}

func (m *UbuntuCephInstallTemplate) Render(pkg urknall.Package) {
	
	ipaddr := GetLocalIP()
	
	pkg.AddCommands("sudoer",
		Shell("echo ' "+Ceph_user+" ALL = (root) NOPASSWD:ALL' | sudo tee /etc/sudoers.d/"+Ceph_user+""),
	)
	pkg.AddCommands("changepermission",
		Shell("sudo chmod 0440 /etc/sudoers.d/"+Ceph_user),
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
		Shell("sudo apt-get -y install ceph-deploy ceph-common ceph-mds dnsmasq openssh-server ntp sshpass "),
	)
	pkg.AddCommands("ipaddress",
		Shell("IP_ADDR='"+ipaddr+"'"),
	)

	pkg.AddCommands("entry",
		Shell("echo 'Adding entry in /etc/hosts'"),
	)
	pkg.AddCommands("edithost",
		Shell("echo '$IP_ADDR "+Host+"'"),
	)
	pkg.AddCommands("ssh",
		Shell("echo 'Processing ssh-keygen'"),
	)

	pkg.AddCommands("sshfile",
		Shell("sudo -u "+Ceph_user+" bash << EOF"+
			"ssh-keygen -N '' -t rsa -f "+User_home+"/.ssh/id_rsa"+
			"cp "+User_home+"/.ssh/id_rsa.pub "+User_home+"/.ssh/authorized_keys"+
			"EOF"),
	)
	pkg.AddCommands("ipKnown_hosts",
		Shell("sudo -H -u "+Ceph_user+" bash -c 'cat > '/"+User_home+"'/.ssh/ssh_config <<EOF'"+
			"'ConnectTimeout 5'"+
			"'Host *''"+
			"'StrictHostKeyChecking no'"+
			"'EOF'"),
	)
	pkg.AddCommands("host",
		Shell("sudo -H -u "+Ceph_user+" bash -c 'cat > '"+User_home+"'/.ssh/config <<EOF'"+
			"'Host "+Host+"'"+
			"'Hostname "+Host+"'"+
			"' User  "+Ceph_user+" '"+
			"'EOF'"),
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
	pkg.AddCommands("osd3",
		Shell("mkdir "+Osd3+"/osd"),
	)
	pkg.AddCommands("getip",
		Shell(" ip3=`echo $IP_ADDR| cut -d'.' -f 1,2,3`"),
	)

	pkg.AddCommands("cephconfig",
		Shell("echo 'Ceph configuration started...'"),
	)
	pkg.AddCommands("conf",
		Shell("sudo -u "+Ceph_user+" bash << EOF"+
			"mkdir "+User_home+"/ceph-cluster"+
			"cd "+User_home+"/ceph-cluster"+

			"ceph-deploy new "+Host+" "+

			"echo 'osd crush chooseleaf type = 0'"+
			"echo 'public network = $ip3.0/24'"+
			"echo 'cluster network = $ip3.0/24'"+
			"ceph-deploy install "+Host+""+
			"ceph-deploy mon create-initial"+
			"ceph-deploy osd prepare "+Host+":"+Osd1+"/osd "+Host+":"+Osd2+"/osd "+Host+":"+Osd3+"/osd"+
			"ceph-deploy osd activate "+Host+":"+Osd1+"/osd "+Host+":"+Osd2+"/osd "+Host+":"+Osd3+"/osd"+
			"ceph-deploy admin "+Host+""+
			"sudo chmod +r /etc/ceph/ceph.client.admin.keyring"+
			"sleep 180"+
			"ceph osd pool set rbd pg_num 150"+
			"sleep 180"+
			"ceph osd pool set rbd pgp_num 150"+
			"EOF"),
	)
	pkg.AddCommands("copy",
		Shell("cp "+User_home+"/ceph-cluster/*.keyring /etc/ceph/"),
	)
	pkg.AddCommands("complete",
		Shell("echo 'Ceph installed successfully.'"),
	)

}


