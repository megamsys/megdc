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
	"github.com/tsuru/config"
	"github.com/astaxie/beego"
	"github.com/megamsys/cloudinabox/controllers"
	"strconv"
)

const (
	version = "0.3.0"
	header  = "Supported-Gulp"
)

const defaultConfigPath = "conf/cib.conf"

func main() {
    beego.SetStaticPath("/static_source", "static_source")
    beego.DirectoryIndex = true
    beego.Router("/", &controllers.LoginController{})
    beego.Router("/signin", &controllers.SignInController{})
    port, _ := config.GetString("beego:http_port")
	if port == "" {
		port = "8085"
	}
	http_port, _ := strconv.Atoi(port)
    beego.HttpPort = http_port
    beego.Run()
}
