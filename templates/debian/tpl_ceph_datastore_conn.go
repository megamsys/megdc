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
	"github.com/megamsys/megdc/templates"
	 u "github.com/megamsys/megdc/templates/ubuntu"
	"github.com/megamsys/urknall"
	"github.com/pborman/uuid"
	"fmt"
)

const (
	Ceph_User = "megdc"
  Poolname = "one"
  Uid =`uuidgen`

Xml=`<secret ephemeral='no' private='no'>
  <uuid>%v</uuid>
  <usage type='ceph'>
          <name>client.libvirt secret</name>
  </usage>
</secret>`
Setval=`sudo virsh secret-set-value --secret %v --base64 $(cat client.libvirt.key)`
Echo =`echo '%v'`
)

var debiancephdatastore *DebianCephDatastore

func init() {
	debiancephdatastore = &DebianCephDatastore{}
	templates.Register("DebianCephDatastore", debiancephdatastore)
}

type DebianCephDatastore struct {}

func (tpl *DebianCephDatastore) Options(t *templates.Template) {}

func (tpl *DebianCephDatastore) Render(p urknall.Package) {
	p.AddTemplate("cephds", &DebianCephDatastoreTemplate{})
}

func (tpl *DebianCephDatastore) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianCephDatastore{})
}

type DebianCephDatastoreTemplate struct {}

func (m *DebianCephDatastoreTemplate) Render(pkg urknall.Package) {
Uid := uuid.NewUUID()
		pkg.AddCommands("cephdatastore",
  	u.AsUser(Ceph_User,u.Shell("ceph osd pool create "+Poolname+" 150")),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool="+Poolname+"'"),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get-key client.libvirt | tee client.libvirt.key"),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get client.libvirt -o ceph.client.libvirt.keyring"),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;cp ceph.client.* /etc/ceph"),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster; "+fmt.Sprintf(Echo,Uid)+" >uid"),
		u.Shell("echo '*****************CEPH_SECRET************************' "),
		u.Shell(fmt.Sprintf(Echo,Uid)),
		u.Shell("echo '*****************************************' "),
		u.WriteFile(UserHomePrefix + Ceph_User + "/ceph-cluster" + "/secret.xml",fmt.Sprintf(Xml,Uid),"root",644),
		u.InstallPackages("libvirt-bin"),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;sudo virsh secret-define secret.xml"),
		u.Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;"+ fmt.Sprintf(Setval,Uid)),
	)

}
