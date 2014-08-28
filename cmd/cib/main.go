/*
** Copyright [2012-2013] [Megam Systems]
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
	"fmt"
	"github.com/megamsys/libgo/cmd"
	"github.com/tsuru/config"
	"log"
	"os"
	"path/filepath"
)

const (
	version = "0.5.0"
	header  = "Supported-CIB"
)

const defaultConfigPath = "conf/cib.conf"

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header)
	m.Register(&GulpStart{m, nil, false}) //start the gulpd daemon
	m.Register(&GulpStop{})               //stop  the gulpd daemon
	m.Register(&GulpUpdate{})             //stop  the gulpd daemon
	return m
}

func main() {
	p, _ := filepath.Abs(defaultConfigPath)
	log.Println(fmt.Errorf("Conf: %s", p))
	config.ReadConfigFile(defaultConfigPath)
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	manager.Run(os.Args[1:])
}


//func main() {
 //   beego.SetStaticPath("../../static_source", "static_source")
  //  beego.DirectoryIndex = true
 //   beego.Router("/", &controllers.LoginController{})
 //   port, _ := config.GetString("beego:http_port")
//	if port == "" {
//		port = "8085"
//	}
//	http_port, _ := strconv.Atoi(port)
  //  beego.HttpPort = http_port
 //   beego.Run()
//}
