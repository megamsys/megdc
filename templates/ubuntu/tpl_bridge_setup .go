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
	BRIDGE_NAME ="bridge"
	PHY_DEV ="phy"
)

var ubuntubridge *UbuntuBridge

func init() {
	ubuntubridge = &UbuntuBridge{}
	templates.Register("UbuntuBridge", ubuntubridge)
}

type UbuntuBridge struct {
	BridgeName string
	PhyDev     string
}

func (tpl *UbuntuBridge) Options(opts map[string]string) {
	if bg, ok := opts[BRIDGE_NAME]; ok {
		tpl.BridgeName = bg
	}

	if ph, ok := opts[PHY_DEV]; ok {
		tpl.BridgeName = ph
	}

}

func (tpl *UbuntuBridge) Render(p urknall.Package) {
	p.AddTemplate("bridge", &UbuntuBridgeTemplate{})
}

func (tpl *UbuntuBridge) Run(target urknall.Target) error {
	return urknall.Run(target, &UbuntuBridge{})
}

type UbuntuBridgeTemplate struct{}

func (m *UbuntuBridgeTemplate) Render(pkg urknall.Package) {

	//    ip := GetLocalIP()

	pkg.AddCommands("setupbrdige",
		Shell(""),
		Shell("sudo echo '"+"%"+"oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-vsctl' >> //etc/sudoers.d/openvswitch"),
		Shell("sudo echo '"+"%"+"oneadmin ALL=(root) NOPASSWD: /usr/bin/ovs-ofctl' >> //etc/sudoers.d/openvswitch"),
		Shell("export BRIDGE_NAME='one'"),
		Shell("export NETWORK_IF='eth0'"),
		Shell("sudo ovs-vsctl add-br one"),
		Shell("sudo echo 'auto one' >> /etc/network/interfaces"),
		Shell("sudo ovs-vsctl add-port one eth0"),

		UpdatePackagesOmitError(),
	)

}
