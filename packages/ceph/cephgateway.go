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

var CEPH_GATEWAY = []string{"CephGateway"}

type Cephgateway struct {
	Fs *gnuflag.FlagSet

	Osd cmd.MapFlag
	CephUser string
	CephPassword string
	IF_name  string
}

func (g *Cephgateway) Info() *cmd.Info {
	desc := `Install Ceph object gateway.

For more information read http://docs.megam.io.`
	return &cmd.Info{
		Name:    "cephgateway",
		Usage:   `cephgateway`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Cephgateway) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "install")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(CEPH_GATEWAY)

	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Cephgateway) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
	}
	return c.Fs
}
