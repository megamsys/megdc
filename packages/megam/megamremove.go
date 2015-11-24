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
	"launchpad.net/gnuflag"
)

var REMOVE_PACKAGES = []string{"SnowflakeRemove",
	"NilavuRemove",
	"GatewayRemove",
	"MegamdRemove"}

type Megamremove struct {
	Fs              *gnuflag.FlagSet
	All             bool
	NilavuRemove    bool
	GatewayRemove   bool
	MegamdRemove    bool
	SnowflakeRemove bool
}

func (g *Megamremove) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "megamremove",
		Usage: `megamremove [--nilavu/-n] [--gateway/-g] [--snowflake/-s]...`,
		Desc: `Remove megam Oja orchestrator. For megdc, available install plaform is ubuntu.
We are working to support centos.
In order to Remove individual packages use the following options.

The [[--nilavu]] parameter removes megam cockpit ui.
This code name is nilavu packaged as megamnilavu.

The [[--gateway]] parameter removes megam gateway apiserver.
This code name is gateway packaged as megamgateway.

The [[--snowflake]] parameter removes megam uidserver.
This code name is snowflake packaged as megamsnowflake.

The [[--megamd]] parameter removes megam omni scheduler.
This code name is megamd packaged as megammegamd.

For more information read http://docs.megam.io.
`,
		MinArgs: 0,
	}
}

func (c *Megamremove) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "removing")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(REMOVE_PACKAGES)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Megamremove) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		/* Remove package commands */
		c.Fs = gnuflag.NewFlagSet("", gnuflag.ExitOnError)
		nilMsg := "Uninstall megam cockpit ui"
		c.Fs.BoolVar(&c.NilavuRemove, "nilavu", false, nilMsg)
		c.Fs.BoolVar(&c.NilavuRemove, "n", false, nilMsg)
		gwyMsg := "Uninstall megam gateway apiserver"
		c.Fs.BoolVar(&c.GatewayRemove, "gateway", false, gwyMsg)
		c.Fs.BoolVar(&c.GatewayRemove, "g", false, gwyMsg)
		megdMsg := "Uninstall megam omni scheduler"
		c.Fs.BoolVar(&c.MegamdRemove, "megamd", false, megdMsg)
		c.Fs.BoolVar(&c.MegamdRemove, "d", false, megdMsg)
		snoMsg := "Uninstall megam uidserver"
		c.Fs.BoolVar(&c.SnowflakeRemove, "snowflake", false, snoMsg)
		c.Fs.BoolVar(&c.SnowflakeRemove, "s", false, snoMsg)
	}
	return c.Fs
}
