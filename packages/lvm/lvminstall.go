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
package lvm
import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
	"launchpad.net/gnuflag"
)

var INSTALL_LVM = []string{"LvmInstall"}

type Lvminstall struct {
	Fs *gnuflag.FlagSet

	Osd cmd.MapFlag
	Bridge string
	PhyDev  string
  Vgname string
}

func (g *Lvminstall) Info() *cmd.Info {
	desc := `Install Lvm with multiple Osds.
  In order to ceph with atleast 1 OSD in a machine use the following options.

  The [--osd osd1=sdb1 --osd osd2=sdb2] parameter describes the storage partitioned all osds. This can be an individual hard disk

  The [[--bride]] parameter defines name of the bride to use for lvm network
  default is br0

  The [[--phy]] parameter defines name of the interface to use for lvm network
  default is eth0

  For more information read http://docs.megam.io.`
  	return &cmd.Info{
  		Name:    "lvminstall",
  		Usage:   `lvminstall`,
  		Desc:    desc,
  		MinArgs: 0,
  	}
  }
  func (c *Lvminstall) Run(context *cmd.Context) error {
  	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "install")
  	w := handler.NewWrap(c)
  	w.IfNoneAddPackages(INSTALL_LVM)

  	if h, err := handler.NewHandler(w); err != nil {
  		return err
  	} else if err := h.Run(); err != nil {
  		return err
  	}
  	return nil
  }

  func (c *Lvminstall) Flags() *gnuflag.FlagSet {
  	if c.Fs == nil {
  		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
  		c.Fs.Var(&c.Osd, "osd", "list of osd storage drive for hosted machine")
  		c.Fs.StringVar(&c.Bridge, "bridge", "br0", "Bridge over Interface")
  		c.Fs.StringVar(&c.Vgname, "vg", "megdc", "password of the ceph user")
      c.Fs.StringVar(&c.PhyDev, "phy", "eth0", "Physical device or Network interface")
  	}
  	return c.Fs
  }
