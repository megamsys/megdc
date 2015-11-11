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

type Cephinstall struct {
	Fs *gnuflag.FlagSet

	Osd1 string
	Osd2 string
}

func (g *Cephinstall) Info() *cmd.Info {
	desc := `Install ceph with 2 OSDs (/storage1 /storage2).

`
	return &cmd.Info{
		Name:    "cephinstall",
		Usage:   `cephinstall`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Cephinstall) Run(context *cmd.Context) error {
	handler.SunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "install")
	w := handler.NewWrap(c)
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
		c.Fs.StringVar(&c.Osd1, "osd1", "", "osd1 storage drive for hosted machine")
		c.Fs.StringVar(&c.Osd2, "osd2", "", "osd2 storage drive for hosted machine")
		//parameter ceph user
		//parameter ceph password
	}
	return c.Fs
}
