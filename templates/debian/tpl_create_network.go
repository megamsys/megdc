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
)

const (
	BRIDGE_NAME = "Bridge"
	PHY_DEV     = "PhyDev"
)

var debiancreatenetwork *DebianCreateNetwork

func init() {
	debiancreatenetwork = &DebianCreateNetwork{}
	templates.Register("DebianCreateNetwork", debiancreatenetwork)
}

type DebianCreateNetwork struct {
	BridgeName string
	PhyDev     string
}

func (tpl *DebianCreateNetwork) Options(t *templates.Template) {
	if bg, ok := t.Options[BRIDGE_NAME]; ok {
		tpl.BridgeName = bg
	}

	if ph, ok := t.Options[PHY_DEV]; ok {
		tpl.PhyDev = ph
	}
}

func (tpl *DebianCreateNetwork) Render(p urknall.Package) {
	p.AddTemplate("createnetwork", &DebianCreateNetworkTemplate{
		BridgeName: tpl.BridgeName,
		PhyDev:     tpl.PhyDev,
	})
}

func (tpl *DebianCreateNetwork) Run(target urknall.Target) error {
	return urknall.Run(target, &DebianCreateNetwork{
		BridgeName: tpl.BridgeName,
		PhyDev:     tpl.PhyDev,
	})
}

type DebianCreateNetworkTemplate struct {
	BridgeName string
	PhyDev     string
}

func (m *DebianCreateNetworkTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("ovs-createnetwork",
		u.Shell("sudo echo '"+"%"+"oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-vsctl' > //etc/sudoers.d/openvswitch"),
		u.Shell("sudo echo '"+"%"+"oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-ofctl' >> //etc/sudoers.d/openvswitch"),
		u.Shell("sudo ovs-vsctl add-br "+m.BridgeName),
		u.Shell("sudo echo 'auto "+m.BridgeName+"' >> /etc/network/interfaces"),
		u.Shell("sudo ovs-vsctl add-port "+m.BridgeName+" "+m.PhyDev+""),
		u.UpdatePackagesOmitError(),
	)

}
