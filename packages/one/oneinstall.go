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

var INSTALL_PACKAGES = []string{"OneInstall"}

type Oneinstall struct {
	Fs         *gnuflag.FlagSet
	OneInstall bool
	Host string
	Username string
	Password string
}

func (g *Oneinstall) Info() *cmd.Info {
	desc := `Install opennebula frontend

Install opennebula frontend (master). This installs opennebula latest release 4.14.
For megdc, available install plaform is ubuntu. We are working to support centos.

For more information read http://docs.megam.io.
`
	return &cmd.Info{
		Name:    "oneinstall",
		Usage:   `oneinstall [--help/-h]...`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Oneinstall) Run(context *cmd.Context) error {
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

func (cmd *Oneinstall) Flags() *gnuflag.FlagSet {
	if cmd.Fs == nil {
		cmd.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		hostMsg := "The host of the server to ssh"
		cmd.Fs.StringVar(&cmd.Host, "host", "localhost", hostMsg)
		usrMsg := "The username of the server"
		cmd.Fs.StringVar(&cmd.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		cmd.Fs.StringVar(&cmd.Password, "password", "", pwdMsg)

	}
	return cmd.Fs
}
