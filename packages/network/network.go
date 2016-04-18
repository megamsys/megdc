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
package network

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"CreateNetworkOpennebula"}

type CreateNetworkOpennebula struct {
	Fs               *gnuflag.FlagSet
	All              bool
CreateNetworkOpennebula   bool
	Bridge string
  Iptype  string
	Ip string
	Size string
	Dns1 string
	Dns2 string
  Gatewayip   string
  Networkmask string
	Username string
	Password string
	Host string
}

func (c *CreateNetworkOpennebula) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "createnetworkopennebula",
		Usage: "createnetworkopennebula",
		Desc: `The createnetworkopennebula is used to setup the network in opennebula.
    In order to install individual packages use the following options.

    The [[--Bridge]] parameter defines to specify your physial bridge.

    The [[--iptype]] parameter defines to specify the ip address ip4 or ip6.
    The [[--ip]] parameter defines to give your ip .
		The [[--size]] parameter defines to specify range of ip address.
	  The [[--dns1]] parameter defines to specify the domain server name.
		The [[--dns2]] parameter defines to specify the domain server name.
		 The [[--gateway]] parameter defines to specify the gateway address.
		 The [[--networkmask]] parameter defines to specify the netmask address.`,
		MinArgs: 0,
	}
}


func (c *CreateNetworkOpennebula) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "installing")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(INSTALL_PACKAGES)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *CreateNetworkOpennebula) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
    bridgeMsg := " To specify physical bridge  in the host where the vm should connect this network interface."
		c.Fs.StringVar(&c.Bridge, "bridge", "one", bridgeMsg)
    iptypeMsg := "To specify ip address ip4 or ip6"
		c.Fs.StringVar(&c.Iptype, "iptype", "IP4", iptypeMsg)
		ipMsg := "To specify the ip"
		c.Fs.StringVar(&c.Ip, "ip", "", ipMsg)
		sizeMsg := "To specify range of ip address."
		c.Fs.StringVar(&c.Size, "size", "", sizeMsg)
		dns1Msg := "To specify domain server address"
		c.Fs.StringVar(&c.Dns1, "dns1", "", dns1Msg)
		dns2Msg := "To specify domain server address"
		c.Fs.StringVar(&c.Dns2, "dns2", "", dns2Msg)
		gatewayMsg := "To specify gateway address"
		c.Fs.StringVar(&c.Gatewayip, "gateway", "", gatewayMsg)
		netmaskMsg := "To specify netmask address"
		c.Fs.StringVar(&c.Networkmask, "networkmask", "", netmaskMsg)
    hostMsg := "To set the hostname"
		c.Fs.StringVar(&c.Host, "host", "localhost", hostMsg)
		usrMsg := "The username of the server"
		c.Fs.StringVar(&c.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		c.Fs.StringVar(&c.Password, "password", "", pwdMsg)
	}

	return c.Fs
}
