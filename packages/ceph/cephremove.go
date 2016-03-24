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

var REMOVE_PACKAGES = []string{"CephRemove"}

type Cephremove struct {
	Fs       *gnuflag.FlagSet
	CephUser string
}

func (g *Cephremove) Info() *cmd.Info {
	desc := `Remove ceph and delete your (2) OSDs.

In order to remove ceph in a machine use the following options.

The [[--cephuser]] parameter defines username chosen as the cephuser. The default user other root can be used.
default is megdc.

For more information read http://docs.megam.io.
`
	return &cmd.Info{
		Name:    "cephremove",
		Usage:   `cephremove`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Cephremove) Run(context *cmd.Context) error {
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

func (c *Cephremove) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		c.Fs.StringVar(&c.CephUser, "cephuser", "megdc", "userid of the ceph user")

	}
	return c.Fs
}
