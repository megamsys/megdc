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
package main

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/megdc/packages/megam"
	//"github.com/megamsys/megdc/packages/config"
//	"github.com/megamsys/megdc/packages/lvm"
	"github.com/megamsys/megdc/packages/one"
	"github.com/megamsys/megdc/packages/onehost"
	"github.com/megamsys/megdc/packages/ceph"
	"github.com/megamsys/megdc/packages/mesos"
	"github.com/megamsys/libgo/cmd"
)

// These variables are populated via the Go linker.
var (
	version string = "1.0"
	commit  string = "01"
	branch  string = "master"
	header  string = "Supported-Megdc"
	date    string
)

func init() {
	// Only log the debug or above
  log.SetLevel(log.DebugLevel)  // level is configurable via cli option.
	// Output to stderr instead of stdout, could also be a file.
  log.SetOutput(os.Stdout)
}

// Only log debug level when the -v flag is passed.
func cmdRegistry(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version+"  "+date, nil, func(modelvl int) {
		if modelvl >= 1 {
			log.SetLevel(log.DebugLevel)
		}
	})
	m.Register(&megam.VerticeInstall{})
	m.Register(&megam.Megamremove{})
	m.Register(&megam.Megamreport{})
	//m.Register(&config.VerticeConf{})
	//m.Register(&lvm.Lvminstall{})
	m.Register(&one.Oneinstall{})
	m.Register(&one.Oneremove{})
	m.Register(&onehost.Onehostinstall{})
	m.Register(&onehost.Onehostremove{})
	m.Register(&ceph.Cephremove{})
	m.Register(&ceph.Cephdatastore{})
	m.Register(&onehost.Createnetwork{})
	m.Register(&onehost.Sshpass{})
	m.Register(&ceph.Cephinstall{})
	m.Register(&ceph.Cephgateway{})
	m.Register(&mesos.MesosMasterInstall{})
	m.Register(&mesos.MesosSlaveInstall{})
	return m
}


//Run the commands from cli.
func main() {
	name := cmd.ExtractProgramName(os.Args[0])
	manager := cmdRegistry(name)
	manager.Run(os.Args[1:])
}
