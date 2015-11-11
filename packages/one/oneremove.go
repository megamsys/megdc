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
package one

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
	"launchpad.net/gnuflag"
)

type Oneremove struct {
	Fs        *gnuflag.FlagSet
	OneRemove bool

}

func (g *Oneremove) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "oneremove",
		Usage:   `oneremove [--help] [--install/-i]...`,
		Desc:  ` Setup Opennebula to your local or remote Machine `,
		MinArgs: 0,
	}
}

func (c *Oneremove) Run(context *cmd.Context) error {
	handler.SunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "remove")
	w := handler.NewWrap(c)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Oneremove) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)

		/* Install package commands */
		oneMsg := "Remove Opennebula"
		c.Fs.BoolVar(&c.OneRemove, "remove", false, oneMsg)
		c.Fs.BoolVar(&c.OneRemove, "r", false, oneMsg)

	}
	return c.Fs
}
