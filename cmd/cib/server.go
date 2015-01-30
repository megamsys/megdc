package main

import (
	"github.com/astaxie/beego"
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/routers/auth"
	"github.com/megamsys/cloudinabox/routers/page"
	"github.com/megamsys/cloudinabox/routers/servers"
	"github.com/megamsys/cloudinabox/routers/nodes"
	"github.com/megamsys/cloudinabox/modules/utils"
//	"github.com/tsuru/config"
	"log"
	"strconv"
	"time"
//	"fmt"
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
	beego.Router("/master", user, "get:Master")
	beego.Router("/dashboard/master", user, "get:MasterDashboard")
	beego.Router("/dashboard/ha", user, "get:HADashboard")
	beego.Router("/dashboard/cs", user, "get:CSDashboard")
	server := new(servers.ServerRouter)
	beego.Router("/servers", server, "get:Get")
	beego.Router("/servers/:id/log", server, "get:Log")
	beego.Router("/servers/getlog", server, "get:GetLog")
	beego.Router("/servers/verify/:name", server, "get:Verify")
	beego.Router("/servers/install/:servername", server, "get:MasterInstall")
	beego.Router("/servers/getIP", server, "get:GetNodeIP")
	beego.Router("/servers/getHAOptions/:ip", server, "get:GetHAOptions")
	beego.Router("/servers/getHA", server, "get:GetHA")
	nodes := new(nodes.NodesRouter)
	beego.Router("/nodes/request/:nodeip", server, "get:NodeInstallRequest")
	beego.Router("/nodes/harequest/:nodeip", server, "get:HANodeInstallRequest")
	beego.Router("/nodes", nodes, "get:Nodes")
	beego.Router("/getnodes", nodes, "get:GetNodes")
	beego.Router("/ha", nodes, "get:Ha")

	//port, _ := config.GetString("beego:http_port")
	port := utils.GetPort()
	if port == "" {
		port = "8077"
	}
	http_port, _ := strconv.Atoi(port)
	beego.HttpPort = http_port
	beego.Run()

}