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
package ceph

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"CephInstall"}

type Cephinstall struct {
	Fs *gnuflag.FlagSet

	Osd cmd.MapFlag
	CephUser string
	CephPassword string
	PhyDev  string
}

func (g *Cephinstall) Info() *cmd.Info {
	desc := `Install ceph with multiple Osds.

In order to ceph with atleast 2 OSDs in a machine use the following options.

The [--osd osd1=storage1 --osd osd2=storage2] parameter describes the storage partitioned all osds. This can be an individual hard disk


The [[--cephuser]] parameter defines username chosen as the cephuser. The default user other root can be used.
default is megdc.

The [[--cephpassword]] parameter defines password for the cephuser.
default is megdc.

The [[--phy]] parameter defines name of the interface to use for ceph network
default is eth0

For more information read http://docs.megam.io.`
	return &cmd.Info{
		Name:    "cephinstall",
		Usage:   `cephinstall`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Cephinstall) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "install")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(INSTALL_PACKAGES)

	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Cephinstall) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		c.Fs.Var(&c.Osd, "osd", "list of osd storage drive for hosted machine")
		c.Fs.StringVar(&c.CephUser, "cephuser", "megdc", "userid used as ceph user")
		c.Fs.StringVar(&c.CephPassword, "cephpassword", "megdc", "password of the ceph user")
    c.Fs.StringVar(&c.PhyDev, "phy", "eth0", "Physical device or Network interface")
	}
	return c.Fs
}
