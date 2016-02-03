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
	"os"

	"github.com/megamsys/megdc/templates"
	u "github.com/megamsys/megdc/templates/ubuntu"
	"github.com/megamsys/urknall"
)

var debiancephremove *DebianCephRemove

func init() {
	debiancephremove = &DebianCephRemove{}
	templates.Register("DebianCephRemove", debiancephremove)
}

type DebianCephRemove struct {
	cephuser string
}

func (tpl *DebianCephRemove) Options(t *templates.Template) {
	if cephuser, ok := t.Options[CephUser]; ok {
		tpl.cephuser = cephuser
	}
}
func (tpl *DebianCephRemove) Render(p urknall.Package) {
	p.AddTemplate("ceph", &DebianCephRemoveTemplate{
        cephuser: tpl.cephuser,
})
}

func (tpl *DebianCephRemove) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianCephRemove{
		cephuser: tpl.cephuser,
	})
}

type DebianCephRemoveTemplate struct {
	cephuser string
}

func (m *DebianCephRemoveTemplate) Render(pkg urknall.Package) {
	host, _ := os.Hostname()

	CephUser := m.cephuser
	pkg.AddCommands("cache-clean",
    u.Shell("rm -r /var/lib/urknall/ceph*"),
	)
	pkg.AddCommands("purgedata",
		u.AsUser(CephUser, u.Shell("ceph-deploy purgedata "+host+"")),
	)
	pkg.AddCommands("forgetKeys",
		u.AsUser(CephUser, u.Shell("ceph-deploy forgetkeys")),
	)
	pkg.AddCommands("purge",
		u.AsUser(CephUser, u.Shell("ceph-deploy purge "+host+"")),
	)
	pkg.AddCommands("rm-sshkey",
		u.AsUser(CephUser, u.Shell("rm -r ~/.ssh")),
	)
	pkg.AddCommands("remove",
		u.Shell("rm -r /var/lib/ceph/"),
		u.Shell("rm -r "+CephUser+"/ceph-cluster"),
		u.Shell("apt-get -y remove ceph-deploy ceph-common ceph-mds"),
		u.Shell("apt-get -y purge ceph-deploy ceph-common ceph-mds"),
		u.Shell("apt-get -y autoremove"),
		u.Shell("rm -r /run/ceph"),
		u.Shell("rm /var/log/upstart/ceph*"),
	)

}
