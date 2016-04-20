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
package attachonehost

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGES = []string{"AttachOneHost"}

type AttachOneHost struct {
	Fs               *gnuflag.FlagSet
	All              bool
  AttachOneHost  bool
	InfoDriver string
	HostName string
  Vm  string
  Networking string
	Username string
	Password string
	Host string
}

func (c *AttachOneHost) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "attachonehost",
		Usage: "attachonehost",
		Desc: `The attachonehost is used to create a host in opennebula.
    In order to install individual packages use the following options.

    The [[--InfoDriver]] parameter defines to set the information driver for the host.
    The [[--Vm]] parameter defines to set the virtual machine manager driver for the host.
    The [[--HostName]] parameter defines to set the hostname.
    The [[--Networking]] parameter defines to set the network driver for the host.`,
		MinArgs: 0,
	}
}


func (c *AttachOneHost) Run(context *cmd.Context) error {
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

func (c *AttachOneHost) Flags() *gnuflag.FlagSet {
	if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
    infodriverMsg := "To set the information driver for the host"
		c.Fs.StringVar(&c.InfoDriver, "infodriver", "kvm", infodriverMsg)
    vmMsg := "To set the virtual machine manager driver for the host"
		c.Fs.StringVar(&c.Vm, "vm", "dummy", vmMsg)
    hostnameMsg := "To set the hostname"
		c.Fs.StringVar(&c.HostName, "hostname", "localhost", hostnameMsg)
    networkingMsg := "To set the network driver for the host"
		c.Fs.StringVar(&c.Networking, "network", "dummy", networkingMsg)
		usrMsg := "The username of the server"
		c.Fs.StringVar(&c.Username, "username", "", usrMsg)
		pwdMsg := "The password of the server"
		c.Fs.StringVar(&c.Password, "password", "", pwdMsg)
		hostMsg := "The host of the server to ssh"
		c.Fs.StringVar(&c.Host, "host", "localhost", hostMsg)
	}

	return c.Fs
}
