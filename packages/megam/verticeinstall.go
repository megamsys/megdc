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
package megam

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"NilavuInstall",
	"GatewayInstall",
	"MegamdInstall"}

type VerticeInstall struct {
	Fs               *gnuflag.FlagSet
	All              bool
	NilavuInstall    bool
	GatewayInstall   bool
	MegamdInstall    bool
	SnowflakeInstall bool
	Host string
	Username string
	Password string
}

func (c *VerticeInstall) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "verticeinstall",
		Usage: "verticeinstall [--nilavu/-n] [--gateway/-g] ",
		Desc: `Install megam Oja orchestrator. For megdc, available install plaform is ubuntu.
We are working to support centos.
In order to install individual packages use the following options.

The [[--nilavu]] parameter defines megam cockpit ui to install.
This code name is nilavu packaged as verticenilavu.

The [[--gateway]] parameter defines megam gateway apiserver to install.
This code name is gateway packaged as verticegateway.

The [[--snowflake]] parameter defines megam uidserver to install.
This code name is snowflake packaged as vertice.

The [[--vertice]] parameter defines megam omni scheduler to install.
This code name is vertice packaged as vertice.

For more information read http://docs.megam.io.`,
		MinArgs: 0,
	}
}

func (c *VerticeInstall) Run(context *cmd.Context) error {
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

func (c *VerticeInstall) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("", gnuflag.ExitOnError)
		nilMsg := "Install megam cockpit ui"
		c.Fs.BoolVar(&c.NilavuInstall, "nilavu", false, nilMsg)
		c.Fs.BoolVar(&c.NilavuInstall, "n", false, nilMsg)
		gwyMsg := "Install megam gateway apiserver"
		c.Fs.BoolVar(&c.GatewayInstall, "gateway", false, gwyMsg)
		c.Fs.BoolVar(&c.GatewayInstall, "g", false, gwyMsg)
		megdMsg := "Install megam omni scheduler"
		c.Fs.BoolVar(&c.MegamdInstall, "vertice", false, megdMsg)
		c.Fs.BoolVar(&c.MegamdInstall, "d", false, megdMsg)
		hostMsg := "The host of the server to ssh"
		c.Fs.StringVar(&c.Host, "host", "localhost", hostMsg)
		usrMsg := "The username of the server"
		c.Fs.StringVar(&c.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		c.Fs.StringVar(&c.Password, "password", "", pwdMsg)
	}
	//	c.Fs = cmd.MergeFlagSet(new(packages.SSHCommand).Flags(), c.Fs)
	return c.Fs
}
