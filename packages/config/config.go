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
package config

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

type VerticeConf struct {
	Fs               *gnuflag.FlagSet
	All				 			 bool
}
func (c *VerticeConf) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "config",
		Usage: "config [-h]",
		Desc: `Configure megam Oja orchestrator. For megdc, available setup.

For more information read http://docs.megam.io.`,
		MinArgs: 0,
	}
}
var CONFIG_PACKAGES= []string{"VerticeConfig"}

func (c *VerticeConf) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "installing")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(CONFIG_PACKAGES)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *VerticeConf) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("", gnuflag.ExitOnError)
	}
	return c.Fs
}
