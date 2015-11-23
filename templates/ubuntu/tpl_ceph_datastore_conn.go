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
	Ceph_User = "megdc"
  Poolname = "one"
  Uid =`uuidgen`
	DsConf = `NAME = "cephds"
DS_MAD = ceph
TM_MAD = ceph
DISK_TYPE = RBD
CEPH_USER = libvirt
CEPH_SECRET = %s
POOL_NAME = %s
BRIDGE_LIST = %s
CEPH_HOST = %s
`
Xml=`<secret ephemeral='no' private='no'>
  <uuid>%s</uuid>
  <usage type='ceph'>
          <name>client.libvirt secret</name>
  </usage>
</secret>
`
)

var ubuntucephdatastore *UbuntuCephDatastore

func init() {
	ubuntucephdatastore = &UbuntuCephDatastore{}
	templates.Register("UbuntuCephDatastore", ubuntucephdatastore)
}

type UbuntuCephDatastore struct {}

func (tpl *UbuntuCephDatastore) Options(opts map[string]string) {}

func (tpl *UbuntuCephDatastore) Render(p urknall.Package) {
	p.AddTemplate("ceph", &UbuntuCephDatastoreTemplate{})
}

func (tpl *UbuntuCephDatastore) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCephDatastore{})
}

type UbuntuCephDatastoreTemplate struct {}

func (m *UbuntuCephDatastoreTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()
	//ip := IP()
		pkg.AddCommands("cephdatastore",
	AsUser(Ceph_User,Shell("ceph osd pool create "+Poolname+" 150")),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool="+Poolname+"'"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get-key client.libvirt | tee client.libvirt.key"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get client.libvirt -o ceph.client.libvirt.keyring"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;cp ceph.client.* /etc/ceph"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;echo "+Uid+" >uid"),
		WriteFile(UserHomePrefix + Ceph_User + "/ceph-cluster" + "/secret.xml",fmt.Sprintf(Xml,Uid),"root",644),
		InstallPackages("libvirt-bin"),
		Shell("sudo virsh secret-define secret.xml"),
		Shell("sudo virsh secret-set-value --secret "+Uid+" --base64 $(cat client.libvirt.key)"),
	)

	pkg.AddCommands("crt-infra",
	  AsUser("oneadmin",Shell("onehost create "+host+" -i kvm -v kvm -n ovswitch")),
		InstallPackages("opennebula-tools"),
		WriteFile("/var/lib/one/ds.conf",fmt.Sprintf(DsConf,Uid,Poolname,host,host),"oneadmin",664),
		AsUser("oneadmin",Shell("onedatastore create /var/lib/one/ds.conf")),
)
}

func (m *UbuntuCephDatastoreTemplate) ip3(a string) string {
	s := strings.Split(IP(""), ".")
	p := s[0 : len(s)-1]
	p = append(p, a)
	return fmt.Sprintf("%s", strings.Join(p, "."))
}
