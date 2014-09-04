package main

import (
	"github.com/astaxie/beego"
	"github.com/megamsys/cloudinabox/routers/servers"
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
	beego.Router("/servernodes/nodes/:nodename", server, "get:NodeInstall")
	//	beego.Router("/servernodes/streamlog", server, "get:StreamLog")


	port, _ := config.GetString("beego:http_port")
	if port == "" {
		port = "8086"
	}
	http_port, _ := strconv.Atoi(port)
	beego.HttpPort = http_port
	beego.Run()

}
