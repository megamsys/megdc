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

package centos

import (
	"github.com/megamsys/megdc/templates"
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

var centoscephdatastore *CentosCephDatastore

func init() {
	centoscephdatastore = &CentosCephDatastore{}
	templates.Register("CentosCephDatastore", centoscephdatastore)
}

type CentosCephDatastore struct {}

func (tpl *CentosCephDatastore) Options(t *templates.Template) {}

func (tpl *CentosCephDatastore) Render(p urknall.Package) {
	p.AddTemplate("cephds", &CentosCephDatastoreTemplate{})
}

func (tpl *CentosCephDatastore) Run(target urknall.Target) error {
	return urknall.Run(target, &CentosCephDatastore{})
}

type CentosCephDatastoreTemplate struct {}

func (m *CentosCephDatastoreTemplate) Render(pkg urknall.Package) {
Uid := uuid.NewUUID()
		pkg.AddCommands("cephdatastore",
  	AsUser(Ceph_User,Shell("ceph osd pool create "+Poolname+" 150")),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get-or-create client.libvirt mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool="+Poolname+"'"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get-key client.libvirt | tee client.libvirt.key"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;ceph auth get client.libvirt -o ceph.client.libvirt.keyring"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;cp ceph.client.* /etc/ceph"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster; "+fmt.Sprintf(Echo,Uid)+" >uid"),
		Shell("echo '*****************************************' "),
		Shell(fmt.Sprintf(Echo,Uid)),
		Shell("echo '*****************************************' "),
		WriteFile(UserHomePrefix + Ceph_User + "/ceph-cluster" + "/secret.xml",fmt.Sprintf(Xml,Uid),"root",644),
		InstallPackages("libvirt-bin"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;sudo virsh secret-define secret.xml"),
		Shell("cd "+UserHomePrefix + Ceph_User+"/ceph-cluster;"+ fmt.Sprintf(Setval,Uid)),
	)
}
