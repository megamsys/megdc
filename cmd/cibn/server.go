package main

import (
	"github.com/astaxie/beego"
	"github.com/megamsys/cloudinabox/routers/servers"
	"github.com/megamsys/cloudinabox/routers/page"
	"github.com/tsuru/config"
	"log"
	"strconv"
	"time"
)


func RunNode() {
	log.Printf("cibnode starting at %s", time.Now())
	handlerNode()
	log.Println("cibnode killed |_|.")
}



func handlerNode() {
	beego.SessionOn = true
	beego.DirectoryIndex = true
	server := new(servers.ServerRouter)
	beego.Router("/servernodes", server, "get:Get")
	beego.Router("/servernodes/:id/log", server, "get:Log")
	beego.Router("/servernodes/getlog", server, "get:GetLog")
	beego.Router("/servernodes/verify/:name", server, "get:Verify")
	beego.Router("/servernodes/nodes/install", server, "get:NodeInstall")
	beego.Router("/servernodes/ha/:name/install", server, "get:HAInstall")
	beego.Router("/servernodes/haproxy/install", server, "post:ProxyInstall")
	beego.Router("/servernodes/devicedetails", server, "get:HADeviceDetails")
	//	beego.Router("/servernodes/streamlog", server, "get:StreamLog")
	user := new(page.PageRouter)
	beego.Router("/dashboard/ha/request", user, "post:HADashboardRequest")
	beego.Router("/dashboard/cs/request", user, "get:CSDashboardRequest")
    beego.Router("/servernodes/nodes/cephoneinstallslave", server, "get:CephOneSlaveInstall")

	port, _ := config.GetString("beego:http_port")
	if port == "" {
		port = "8078"
	}
	http_port, _ := strconv.Atoi(port)
	beego.HttpPort = http_port
	beego.Run()

}
