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
	"github.com/megamsys/cloudinabox/routers/auth"
	"github.com/megamsys/cloudinabox/routers/page"
	"github.com/megamsys/cloudinabox/routers/servers"
	"github.com/megamsys/cloudinabox/models/orm"
//	"github.com/megamsys/cloudinabox/setting"
	"strconv"
)

const (
	version = "0.3.0"
	header  = "Supported-Gulp"
)

const defaultConfigPath = "conf/cib.conf"

func main() {
	//setting.LoadConfig()
	
	// set db 	
    db := orm.OpenDB()
    dbmap := orm.GetDBMap(db)
    // initialize the DbMap
    orm.InitDB(dbmap)
    defer db.Close()
    
	beego.SessionOn = true
    beego.SetStaticPath("/static_source", "static_source")
    beego.DirectoryIndex = true
    login := new(auth.LoginRouter)
	beego.Router("/", login, "get:Get")
	beego.Router("/signin", login, "post:Login")
	beego.Router("/logout", login, "get:Logout")
	user := new(page.PageRouter)
	beego.Router("/index", user, "get:Get")
	server := new(servers.ServerRouter)
	beego.Router("/servers", server, "get:Get")
	beego.Router("/servers/1", server, "get:Progress")
	beego.Router("/servers/join", server, "get:Join")
	
    port, _ := config.GetString("beego:http_port")
	if port == "" {
		port = "8085"
	}
	http_port, _ := strconv.Atoi(port)
    beego.HttpPort = http_port
    beego.Run()
}
