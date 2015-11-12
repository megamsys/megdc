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

type Oneinstall struct {
	Fs         *gnuflag.FlagSet
	OneInstall bool
}

func (g *Oneinstall) Info() *cmd.Info {
	desc := `Install opennebula frontend
`
	return &cmd.Info{
		Name:    "oneinstall",
		Usage:   `oneinstall [--help/-h]...`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Oneinstall) Run(context *cmd.Context) error {
	handler.SunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "install")
	w := handler.NewWrap(c)
	c.oneInstall(w)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Oneinstall) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		/* Install package commands
		oneMsg := "Install Opennebula"
		c.Fs.BoolVar(&c.OneInstall, "install", false, oneMsg)
		c.Fs.BoolVar(&c.OneInstall, "i", false, oneMsg)*/
	}
	return c.Fs
}
func (c *Oneinstall) oneInstall(w *handler.WrappedParms) {
	DEFAULT_PACKAGES := []string{"OneInstall"}

	if w.Empty() {
		for i := range DEFAULT_PACKAGES {
			w.AddPackage(DEFAULT_PACKAGES[i])
		}
	}
}
