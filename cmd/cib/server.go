package main

import (
	"github.com/astaxie/beego"
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/routers/auth"
	"github.com/megamsys/cloudinabox/routers/page"
	"github.com/megamsys/cloudinabox/routers/servers"
	"github.com/tsuru/config"
	"log"
	"strconv"
	"time"
)


func RunWeb() {
	log.Printf("cib starting at %s", time.Now())
	handlerWeb()
	log.Println("cib killed |_|.")
}



func handlerWeb() {
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
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
	beego.Router("/dash", user, "get:Dash")
	server := new(servers.ServerRouter)
	beego.Router("/servers", server, "get:Get")
	beego.Router("/servers/:id/log", server, "get:Log")
	beego.Router("/servers/getlog", server, "get:GetLog")
	beego.Router("/servers/verify/:name", server, "get:Verify")
	beego.Router("/servers/:servername", server, "get:MasterInstall")
	beego.Router("/servers/nodes/:nodename", server, "get:NodesInstall")

	port, _ := config.GetString("beego:http_port")
	if port == "" {
		port = "8086"
	}
	http_port, _ := strconv.Atoi(port)
	beego.HttpPort = http_port
	beego.Run()

}