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

package servers

import (
//	"strings"
//	"github.com/megamsys/cloudinabox/modules/utils"
    "github.com/astaxie/beego"
    "github.com/megamsys/cloudinabox/routers/base"
    "github.com/megamsys/cloudinabox/modules/servers"
    "github.com/megamsys/cloudinabox/models/orm"
    "github.com/megamsys/cloudinabox/models"
    "github.com/ActiveState/tail"
    "github.com/gorilla/websocket"
     "net/http"
     "container/list"
     "encoding/json"
	"time"
	"strconv"
 //   "regexp"
   "fmt"
)

var serversList = [...]string{ "MEGAM", "COBBLER"}

// PageRouter serves home page.
type ServerRouter struct {
	base.BaseRouter
}

// Get implemented dashboard page.
func (this *ServerRouter) Get() {
	var servers orm.Servers
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	this.Data["IsLoginPage"] = true
	this.Data["Username"] = "Megam"
	this.TplNames = "servers/servers.html" 
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	n := len(serversList)
	for i := 0; i < n; i++ {
       err := dbmap.SelectOne(&servers, "select * from servers where Name=?", serversList[i])	  
	   fmt.Println(err)
    }
	result["success"] = true
	result["data"] = servers
	
}

func (this *ServerRouter) Log() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "servers/log.html" 
	server := this.Ctx.Input.Param(":id")
	this.Data["ServerName"] = server
}

func (this *ServerRouter) GetLog() {
	uname := this.GetString("uname")
	fmt.Println("Join entry")
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

// Join method handles WebSocket requests for WebSocketController.
func (this *ServerRouter) Join() {
	uname := this.GetString("uname")
	fmt.Println("Join entry")
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
	startPolling()
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
	fmt.Println("broadcast entry")
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
var i = 0
var oldServerName = ""
func doSomething(server string) bool { 
	    if oldServerName == "" {
	    	oldServerName = server
	    	i = i + 3
	    	PublishMessage(server, i)
	    } else {
	    	if oldServerName == server {
	    		 i = i + 3
	             if (i <= 99) {
	                PublishMessage(server, i)
                 } else {
                 	return false
                 }
	             //else {
                 //	publish <- newEvent(models.EVENT_MESSAGE, "megam", "completed")
                 // } 
	    	} else {
	    		oldServerName = server
	    		i = 0
	    		i = i + 3
	    		PublishMessage(server, i)
	    	}
	    }
	    return true
} 

func publishLog(server string) {
	t, _ := tail.TailFile("/var/log/opennebula.log", tail.Config{Follow: true})
        for line := range t.Lines {
          fmt.Println(line.Text)
          publish <- newEvent(models.EVENT_MESSAGE, server, line.Text)
        //  t.Stop()
        }
}

func PublishMessage(server string, i int) {
	 publish <- newEvent(models.EVENT_MESSAGE, server, strconv.Itoa(i))
     fmt.Printf("doing something: %v", i)
}

func startPolling() { 
	  //  i := 0
/*	  for i := range serversList {
        for _ = range time.Tick(1 * time.Second) { 
               if doSomething(serversList[i]) == false {
               	  break
               } 
         } 
        }*/
	    db := orm.OpenDB()
	    dbmap := orm.GetDBMap(db)
	    var server orm.Servers
	    fmt.Println("polling entry")
	    for i := range serversList {
	    	fmt.Println(serversList[i])
            err := dbmap.SelectOne(&server, "select * from servers where Name=?", serversList[i])	  
	        fmt.Println(err)
	        if server.Install != true {
	   	       err := servers.InstallServers(serversList[i])
	   	       if err != nil {
	   		
	   	       }
	        }
         }
	    
//	    t, _ := tail.TailFile("/var/log/syslog", tail.Config{Follow: true})
 //       for line := range t.Lines {
 //         fmt.Println(line.Text)
  //        t.Stop()
  //      }
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
				//oldUserReconnect(subscribers, sub.Name)
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

/*func oldUserReconnect(subscribers *list.List, user string) bool {
	  for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == user {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub, "") // Publish a LEAVE event.
					return true
				}
			}
	  return false
}*/

