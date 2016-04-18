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
package datastore

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"CreateDatastoreLvm"}

type CreateDatastoreLvm struct {
	Fs               *gnuflag.FlagSet
	All              bool
	CreateDatastoreLvm    bool
	PoolName string
  VgName  string
	Hostname string
	Username string
	Password string
	Host string
}

func (c *CreateDatastoreLvm) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "createdatastorelvm",
		Usage: "createdatastorelvm",
		Desc: `The createdatestorelvm is used to create a datastore for lvm.
    In order to install individual packages use the following options.

    The [[--poolname]] parameter defines to specify the name storing data object.

    The [[--vgname]] parameter defines that LVM volume group name.
    The [[--Hostname]] parameter defines to specify the hostname.`,
		MinArgs: 0,
	}
}


func (c *CreateDatastoreLvm) Run(context *cmd.Context) error {
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

func (c *CreateDatastoreLvm) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
    poolnameMsg := " specify name storing data object"
		c.Fs.StringVar(&c.PoolName, "poolname", "one", poolnameMsg)
    vgMsg := "To specify LVM volume group name"
		c.Fs.StringVar(&c.VgName, "vgname", "vg-one-0", vgMsg)
    hostnameMsg := "To set the hostname"
		c.Fs.StringVar(&c.Hostname, "hostname", "localhost", hostnameMsg)
		usrMsg := "The username of the server"
		c.Fs.StringVar(&c.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		c.Fs.StringVar(&c.Password, "password", "", pwdMsg)
		hostMsg := "The host of the server to ssh"
		c.Fs.StringVar(&c.Host, "host", "localhost", hostMsg)
	}

	return c.Fs
}
