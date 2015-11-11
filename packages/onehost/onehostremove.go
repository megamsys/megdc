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
package onehost

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
	"launchpad.net/gnuflag"
)

type Onehostremove struct {
	Fs *gnuflag.FlagSet

	Host     string
	Username string
	Password string
}

func (g *Onehostremove) Info() *cmd.Info {
	desc := `Remove Opennebula host.
`
	return &cmd.Info{
		Name:    "onehostremove",
		Usage:   `onehostremove [--host] [--username]...`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Onehostremove) Run(context *cmd.Context) error {
	handler.SunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "remove")
	w := handler.NewWrap(c)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Onehostremove) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)

		c.Fs.StringVar(&c.Host, "host", "", "host address for machine")
		c.Fs.StringVar(&c.Host, "h", "", "host address for machine")
		c.Fs.StringVar(&c.Username, "username", "", "username for hosted machine")
		c.Fs.StringVar(&c.Username, "u", "", "username for hosted machine")
		c.Fs.StringVar(&c.Password, "password", "", "password for hosted machine")
		c.Fs.StringVar(&c.Password, "p", "", "password for hosted machine")
	}
	return c.Fs
}
