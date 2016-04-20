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
package bridge

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"CreateBridge"}

type CreateBridge struct {
	Fs               *gnuflag.FlagSet
	All              bool
	CreateBridge    bool
  Bridgename string
  PhyDev string
	Network string
	Netmask string
	Gateway string
	Dnsname1 string
	Dnsname2 string
	Host string
	Username string
	Password string
  }

func (c *CreateBridge) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "createbridge",
		Usage: "createbridge",
		Desc: ` Create bridge and interfaces.
    In order to install individual packages use the following options.
    The [[--bridgename]] parameter defines to specify which name the bridge is created.
    The [[--phydev]] parameter defines to specify the interface name
		The [[--network]] parameter defines to specify the network address
		The [[--netmask]] parameter defines to specify the netmask address
		The [[--gateway]] parameter defines to specify the gateway address
    The [[--dnsname1]] parameter defines to specify the  name server
		The [[--dnsname2]] parameter defines to specify the  name server`,
		MinArgs: 0,
	}
}


func (c *CreateBridge) Run(context *cmd.Context) error {
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

func (c *CreateBridge) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
    bridgeMsg := "specify the name of bridge"
		c.Fs.StringVar(&c.Bridgename, "bridgename", "one", bridgeMsg)
    phydevMsg := "specify the interfacename"
		c.Fs.StringVar(&c.PhyDev, "phydev", "eth0", phydevMsg)
		networkMsg := "specify the network address"
		c.Fs.StringVar(&c.Network, "network", "", networkMsg)
		netmaskMsg := "specify the netmask address"
		c.Fs.StringVar(&c.Netmask, "netmask", "", netmaskMsg)
		gatewayMsg := "specify the gateway address"
		c.Fs.StringVar(&c.Gateway, "gateway", "", gatewayMsg)
		dnsname1Msg := "specify the name server"
		c.Fs.StringVar(&c.Dnsname1, "dnsname1", "", dnsname1Msg)
		dnsname2Msg := "specify the name server"
		c.Fs.StringVar(&c.Dnsname2, "dnsname2", "", dnsname2Msg)
		hostMsg := "The host of the server to ssh"
		c.Fs.StringVar(&c.Host, "host", "localhost", hostMsg)
		usrMsg := "The username of the server"
		c.Fs.StringVar(&c.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		c.Fs.StringVar(&c.Password, "password", "", pwdMsg)
	}

	return c.Fs
}
