/*
** Copyright [2012-2014] [Megam Systems]
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

package servers

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/ActiveState/tail"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/megamsys/cloudinabox/models"
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/modules/servers"
	"github.com/megamsys/cloudinabox/routers/base"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var serversList = [...]string{"MEGAM", "COBBLER", "OPENNEBULA", "OPENNEBULAHOST"}
const layout = "Jan 2, 2006 at 3:04pm (MST)"

// PageRouter serves home page.
type ServerRouter struct {
	base.BaseRouter
}

// Get implemented dashboard page.
func (this *ServerRouter) Get() {
	var servers orm.Servers
	var servers_output []string

	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

	this.Data["IsLoginPage"] = true
	this.Data["Username"] = this.GetUser()
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	n := len(serversList)
	for i := 0; i < n; i++ {
		log.Printf("[%s] Selecting\n", serversList[i])
		err := dbmap.SelectOne(&servers, "select * from servers where Name=?", serversList[i])
		if err != nil {
			tmpserver := &orm.Servers{0, serversList[i], false, "", "", ""}
			jsonMsg, _ := json.Marshal(tmpserver)
			servers_output = append(servers_output, string(jsonMsg))
			log.Printf("[%s] SQL select  {%s}\n%v\n", serversList[i], jsonMsg, servers_output)
		} else {
			jsonMsg, _ := json.Marshal(servers)
			servers_output = append(servers_output, string(jsonMsg))
			log.Printf("[%s] SQL select ignored {%s}\n%v\n", serversList[i], jsonMsg, servers_output)
			
		}
	}
	result["success"] = true
	result["data"] = servers_output

}

func (this *ServerRouter) MasterInstall() {
	var server orm.Servers
	var servers_output string
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

	this.Data["IsLoginPage"] = true
	this.Data["Username"] = this.GetUser()
	servername := this.Ctx.Input.Param(":servername")

	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	err := dbmap.SelectOne(&server, "select * from servers where Name=?", servername)
	fmt.Println(err)
	if server.Install != true {
		//if len(server.IP) > 0 {
		//	node_err := servers.InstallNode(&server)
		//	fmt.Printf("%s", node_err)
		//	if node_err != nil {
		//		result["success"] = false
		//	} else {
				result["success"] = true
		//	}
		//} else {
			err := servers.InstallServers(servername)
			fmt.Printf("%s", err)
			if err != nil {
				result["success"] = false
			} else {
				result["success"] = true
			}
	//	}
	} else {
		result["success"] = true
	}
	dberr := dbmap.SelectOne(&server, "select * from servers where Name=?", servername)
	if dberr != nil {
			tmpserver := &orm.Servers{0, servername, false, "", "", ""}
			jsonMsg, _ := json.Marshal(tmpserver)
			servers_output = string(jsonMsg)
			log.Printf("[%s] SQL select  {%s}\n%v\n", servername, jsonMsg, servers_output)
	} else {
			jsonMsg, _ := json.Marshal(server)
			servers_output = string(jsonMsg)
			log.Printf("[%s] SQL select ignored {%s}\n%v\n", servername, jsonMsg, servers_output)
	}
	result["data"] = servers_output
}


func (this *ServerRouter) NodeInstallRequest() {
	result := map[string]interface{}{
		"success": false,
	}
    nodeip := this.Ctx.Input.Param(":nodeip")
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	var nodes orm.Nodes
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	node_err := servers.InstallNode(nodeip)
	fmt.Printf("%s", node_err)
	if node_err != nil {
		result["success"] = false
	} else {
		uerr := updateNode(nodeip)
		fmt.Println(uerr)
		err := dbmap.SelectOne(&nodes, "select * from nodes where IP=?", nodeip)
		if err == nil {
			jsonMsg, _ := json.Marshal(nodes)
			result["data"] = string(jsonMsg)
		}
		result["success"] = true
	}
}


func (this *ServerRouter) NodeInstall() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	//servername := this.Ctx.Input.Param(":nodename")
	
    servername := "NODEINSTALL"
	err := servers.InstallServers(servername)
	fmt.Printf("%s", err)
	if err != nil {
		result["success"] = false
	} else {
		result["success"] = true
	}

}

func (this *ServerRouter) Log() {
	fmt.Println("Join entry LOG()============> ")
	fmt.Println(this.Ctx.Input.Param(":id"))
	this.Data["IsLoginPage"] = true
	this.TplNames = "servers/log.html"
	server := this.Ctx.Input.Param(":id")
	this.Data["ServerName"] = server
}

func (this *ServerRouter) GetLog() {
	uname := this.GetString("uname")
	fmt.Println("Join entry ============> ")
	fmt.Println(uname)
	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	// Join chat room.
	Join(uname, ws)
	publishLog(uname)
	defer Leave(uname)

}

func (this *ServerRouter) Verify() {
	servername := this.Ctx.Input.Param(":name")
	fmt.Println(servername)
	var server orm.Servers
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	err := dbmap.SelectOne(&server, "select * from servers where Name=?", servername)

	if err != nil {
		result["failure_message"] = err.Error()
	}
	fmt.Println(server.Install)
	if !server.Install {
		result["success"] = false
	} else {
		result["success"] = true
	}
}

func (this *ServerRouter) GetNodeIP() {
	var node orm.Nodes
	//nodename := this.Ctx.Input.Param(":nodename")
	filePath := "/var/lib/megam/megamcib/boxips"

	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)

	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("no such file or directory: %s", filePath)
			// open output file
			_, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
		}

		t, err := tail.TailFile(filePath, tail.Config{Follow: true, MustExist: true})
		if err != nil {
			log.Printf("ERROR LOG READ ==> : %s", err.Error())
		}

		for line := range t.Lines {
			fmt.Println("-------------------------")
			fmt.Println(line.Text)
			err1 := dbmap.SelectOne(&node, "select * from nodes where IP=?", line.Text)
			fmt.Println("++++++++++++++++++++++++")
			fmt.Println(err1)
			if err1 != nil {
				newnode := orm.NewNode(line.Text)
				orm.ConnectToTable(dbmap, "nodes", newnode)
				err := dbmap.Insert(&newnode)
				if err != nil {
					fmt.Println("Node insert error======>")
					result["ip"] = false
					result["ipvalue"] = ""
				}
				result["ip"] = true
				result["ipvalue"] = line.Text
			}
			if len(line.Text) > 0 {
				err := os.Remove(filePath)

				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
}

/*func verify_NODES(nodename string) bool {
	var node orm.Nodes
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	err := dbmap.SelectOne(&server, "select * from servers where Name=?", nodename)

	if err != nil {
		return false
	} else {
		return true
	}
}*/

// Join method handles WebSocket requests for WebSocketController.
func (this *ServerRouter) Join() {
	uname := this.GetString("uname")
	fmt.Println("Join entry ===================>  ")
	fmt.Println(uname)
	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	// Join chat room.
	Join(uname, ws)
	//	startPolling()
	defer Leave(uname)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	fmt.Println("broadcast entry ========> ")
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}


func publishLog(server string) {
	fmt.Printf("LOG FILE NAME ===========> : %s", server)
	t, err := tail.TailFile("/var/log/megam/megamcib/"+server+".log", tail.Config{Follow: true})
	if err != nil {
		log.Printf("ERROR LOG READ ==> : %s", err.Error())
	}
	for line := range t.Lines {
		fmt.Println(line.Text)
		publish <- newEvent(models.EVENT_MESSAGE, server, line.Text)
		//  t.Stop()
	}
}

func PublishMessage(server string, i int) {
	publish <- newEvent(models.EVENT_MESSAGE, server, strconv.Itoa(i))
}


type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

func newEvent(ep models.EventType, user, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.
	unsubscribe = make(chan string, 10)
	// Send events here to publish them.
	publish = make(chan models.Event, 10)
	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()
)

func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:
			// Notify waiting list.
			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.Server, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub, "") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}

func updateNode(nodeip string) error {
        // insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)

		node := orm.Nodes{}
		err := dbmap.SelectOne(&node, "select * from nodes where IP=?", nodeip)
		if err != nil {
			fmt.Println("node select error======>")
			return err
		}
		err3 := orm.DeleteRowFromNodeIP(dbmap, nodeip)
		if err3 != nil {
			log.Printf("node delete error")
			return err3
		}
		time := time.Now()
		update_node := orm.Nodes{Id: node.Id, Install: true, IP: node.IP, InstallDate: node.InstallDate, UpdateDate: time.Format(layout)}
		orm.ConnectToTable(dbmap, "nodes", update_node)
		err2 := dbmap.Insert(&update_node)
	
		if err2 != nil {
			fmt.Println("node insert error======>")
			return err2
		}	
	return nil
}	

