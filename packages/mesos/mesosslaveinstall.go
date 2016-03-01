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
package mesos

import (
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/megdc/handler"
//	"github.com/megamsys/megdc/packages"
	"launchpad.net/gnuflag"
)

var INSTALL_PACKAGE = []string{"MesosSlaveInstall"}

type MesosSlaveInstall struct {
	Fs       *gnuflag.FlagSet
	MasterIp string
  SparkVersion string
}

func (c *MesosSlaveInstall) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "mesosslaveinstall",
		Usage:   `mesosslaveinstall`,
		Desc: `Install mesosmaster for cluster manager that provides efficient resource isolation and sharing
    across distributed applications or framework.
In order to install messosslave to use the following options.

The [[--masterip]] parameter defines ip for the mesosmaster.
The [[--sparkversion]] parameter defines version for the spark.
For more information read http://docs.megam.io.`,

  MinArgs: 0,
}
	}


func (c *MesosSlaveInstall) Run(context *cmd.Context) error {
	handler.FunSpin(cmd.Colorfy(handler.Logo, "green", "", "bold"), "", "installing")
	w := handler.NewWrap(c)
	w.IfNoneAddPackages(INSTALL_PACKAGE)
	if h, err := handler.NewHandler(w); err != nil {
		return err
	} else if err := h.Run(); err != nil {
		return err
	}
	return nil
}

func (c *MesosSlaveInstall) Flags() *gnuflag.FlagSet {
  if c.Fs == nil {
		c.Fs = gnuflag.NewFlagSet("megdc", gnuflag.ExitOnError)
		c.Fs.StringVar(&c.MasterIp, "masterip", "", "ip of the mesosmaster")
    	c.Fs.StringVar(&c.SparkVersion, "sparkversion", "1.5.1", "version of the spark")
	}
	return c.Fs
}
