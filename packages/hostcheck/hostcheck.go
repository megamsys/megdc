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
package hostcheck

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"HostCheck"}

type HostCheck struct {
	Fs               *gnuflag.FlagSet
	All              bool
	HostCheck    bool
	Host string
	Username string
	Password string
}

func (c *HostCheck) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "hostcheck",
		Usage: "hostcheck",
		Desc: `The hostcheck is used to check that vt enable`,
		MinArgs: 0,
	}
}


func (c *HostCheck) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "installing")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(INSTALL_PACKAGES)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *HostCheck) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		hostMsg := "The host of the server to ssh"
		c.Fs.StringVar(&c.Host, "host", "localhost", hostMsg)
		usrMsg := "The username of the server"
		c.Fs.StringVar(&c.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		c.Fs.StringVar(&c.Password, "password", "", pwdMsg)
	}

	return c.Fs
}
