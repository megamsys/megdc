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
	"github.com/megamsys/megdc/templates"
	"github.com/megamsys/urknall"
)

const (
	BRIDGE_NAME = "Bridge"
	PHY_DEV     = "PhyDev"
)

var ubuntucreatenetwork *UbuntuCreateNetwork

func init() {
	ubuntucreatenetwork = &UbuntuCreateNetwork{}
	templates.Register("UbuntuCreateNetwork", ubuntucreatenetwork)
}

type UbuntuCreateNetwork struct {
	BridgeName string
	PhyDev     string
}

func (tpl *UbuntuCreateNetwork) Options(opts map[string]string) {
	if bg, ok := opts[BRIDGE_NAME]; ok {
		tpl.BridgeName = bg
	}

	if ph, ok := opts[PHY_DEV]; ok {
		tpl.PhyDev = ph
	}
}

func (tpl *UbuntuCreateNetwork) Render(p urknall.Package) {
	p.AddTemplate("createnetwork", &UbuntuCreateNetworkTemplate{
		BridgeName: tpl.BridgeName,
		PhyDev:     tpl.PhyDev,
	})
}

func (tpl *UbuntuCreateNetwork) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuCreateNetwork{
		BridgeName: tpl.BridgeName,
		PhyDev:     tpl.PhyDev,
	})
}

type UbuntuCreateNetworkTemplate struct {
	BridgeName string
	PhyDev     string
}

func (m *UbuntuCreateNetworkTemplate) Render(pkg urknall.Package) {
	pkg.AddCommands("ovs-createnetwork",
		Shell("sudo echo '"+"%"+"oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-vsctl' > //etc/sudoers.d/openvswitch"),
		Shell("sudo echo '"+"%"+"oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-ofctl' >> //etc/sudoers.d/openvswitch"),
		Shell("sudo ovs-vsctl add-br "+m.BridgeName),
		Shell("sudo echo 'auto "+m.BridgeName+"' >> /etc/network/interfaces"),
		Shell("sudo ovs-vsctl add-port "+m.BridgeName+" "+m.PhyDev+""),
		UpdatePackagesOmitError(),
	)

}
