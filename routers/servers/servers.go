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
	"bytes"
	"github.com/ActiveState/tail"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/megamsys/cloudinabox/models"
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/modules/servers"
	"github.com/megamsys/cloudinabox/routers/base"
	"github.com/megamsys/libgo/exec"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"time"
	"io/ioutil"
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
		err := dbmap.SelectOne(&servers, "select * from servers where Stype='MASTER' and Name=?", serversList[i])
		if err != nil {
			tmpserver := &orm.Servers{0, serversList[i], false, "", "", "", "", ""}
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

// Get implemented dashboard page.
func (this *ServerRouter) GetHA() {
	var serverlist []orm.Servers
	var servers_output []orm.Servers

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
	//n := len(serversList)
	_, err := dbmap.Select(&serverlist, "select * from servers where Stype='HA' GROUP BY IP") 
	if err != nil {
		result["success"] = false
		result["data"] = servers_output
	} else {
		result["success"] = true
		result["data"] = serverlist
  	}
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
	err := dbmap.SelectOne(&server, "select * from servers where Stype='MASTER' and Name=?", servername)
	fmt.Println(err)
	if server.Install != true {
			err := servers.InstallServers(servername)
			fmt.Printf("%s", err)
			if err != nil {
				result["success"] = false
			} else {
				newserver := orm.NewServer(servername, "localhost", "MASTER", "")
				orm.ConnectToTable(dbmap, "servers", newserver)
				derr := dbmap.Insert(&newserver)

				if derr != nil {
					fmt.Println("server insert error======>")
				}
				uerr := updateServer("localhost", servername)
				if uerr != nil {
					fmt.Println("server insert error======>")
					result["success"] = false
				}
				
				haerr := servers.CheckHAInstall(servername)
				if haerr != nil {
					result["success"] = false
				}
				
				result["success"] = true
			}
	} else {
		result["success"] = true
	}
	dberr := dbmap.SelectOne(&server, "select * from servers where Stype='MASTER' and Name=?", servername)
	if dberr != nil {
			tmpserver := &orm.Servers{0, servername, false, "", "", "", "", ""}
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
	node_err := servers.InstallNode(nodeip, "COMPUTE", "")
	fmt.Printf("%s", node_err)
	if node_err != nil {
		result["success"] = false
	} else {
		    nodename := "" 
			newnode := orm.NewNode(nodeip, nodename)
			orm.ConnectToTable(dbmap, "nodes", newnode)
			err := dbmap.Insert(&newnode)
			if err != nil {
				fmt.Println("Node insert error======>")
				jsonMsg, _ := json.Marshal(nodes)
				result["data"] = string(jsonMsg)
			}
			//uerr := updateNode(nodeip)
			//fmt.Println(uerr)
			nerr := dbmap.SelectOne(&nodes, "select * from nodes where IP=?", nodeip)
			if nerr == nil {
				jsonMsg, _ := json.Marshal(nodes)
				result["data"] = string(jsonMsg)
			}
		result["success"] = true
	}
}


func (this *ServerRouter) HANodeInstallRequest() {
	result := map[string]interface{}{
		"success": false,
	}
	nodeip := this.Ctx.Input.Param(":nodeip")
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	var serverlist []orm.Servers
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	
		_, err := dbmap.Select(&serverlist, "select * from servers where Stype='MASTER'")      
        
		if err != nil {
			result["success"] = false
		} else {
			nodename := ""
			for _, p := range serverlist {
				newserver := orm.NewServer(p.Name, nodeip, "HA", nodename)
				orm.ConnectToTable(dbmap, "servers", newserver)
				derr := dbmap.Insert(&newserver)
				if derr != nil {
					fmt.Println("server insert error======>")
				}
    		}
		 	for x, p := range serverlist {
          		log.Printf("    %d: %v\n", x, p)
          		node_err := servers.InstallNode(nodeip, "HA", p.Name)
				fmt.Printf("%s", node_err)
				if node_err != nil {
					result["success"] = false
				} else {
					uerr := updateServer(nodeip, p.Name)
				    if uerr != nil {
						fmt.Println("server insert error======>")
						result["success"] = false
					}
					result["success"] = true
					}
    		}
		 	var haserver orm.HAServers
		 	nerr := dbmap.SelectOne(&haserver, "select * from haservers where nodeip2=?", nodeip)
			if nerr != nil {
				result["success"] = false
			} else {
				masterproxyerr := servers.InstallProxy(&haserver, "MASTER")
				if masterproxyerr != nil {
					result["success"] = false
				} 
				
				jsonMsg, _ := json.Marshal(haserver)
				
		 		url := "http://" + nodeip + ":8078/servernodes/haproxy/install"
		 		res, rerr := http.Post(url, "string", bytes.NewBufferString(string(jsonMsg)))		
	    		if rerr != nil {
		  	 	 	result["success"] = false
	    		} else {
		 			if res.StatusCode > 299 {
						result["success"] = false
					} else {
						result["success"] = true
					}
	  			}
	    	}
	}
}

/*func (this *ServerRouter) NodeInstallRequest() {
	result := map[string]interface{}{
		"success": false,
	}
	req := this.Ctx.Input.Param(":nodeip")
	s := strings.Split(req, "-")
    nodedata := s[1]
    nodetype := s[0]
    ns := strings.Split(nodedata, "=")
    nodeip := ns[1]
    nodename := ns[0]
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	var nodes orm.Nodes
	var serverlist []orm.Servers
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	if nodetype == "COMPUTE" {
		node_err := servers.InstallNode(nodeip, nodetype, "")
		fmt.Printf("%s", node_err)
		if node_err != nil {
			result["success"] = false
		} else {
			if nodetype == "COMPUTE" {
				newnode := orm.NewNode(nodeip, nodename)
				orm.ConnectToTable(dbmap, "nodes", newnode)
				err := dbmap.Insert(&newnode)
				if err != nil {
					fmt.Println("Node insert error======>")
					jsonMsg, _ := json.Marshal(nodes)
					result["data"] = string(jsonMsg)
				}
				//uerr := updateNode(nodeip)
				//fmt.Println(uerr)
				nerr := dbmap.SelectOne(&nodes, "select * from nodes where IP=?", nodeip)
				if nerr == nil {
					jsonMsg, _ := json.Marshal(nodes)
					result["data"] = string(jsonMsg)
				}
			} 
			result["success"] = true
		}
	} else {
		_, err := dbmap.Select(&serverlist, "select * from servers where Stype='MASTER'")      
        
		if err != nil {
			result["success"] = false
		} else {
			for _, p := range serverlist {
				newserver := orm.NewServer(p.Name, nodeip, "HA", nodename)
				orm.ConnectToTable(dbmap, "servers", newserver)
				derr := dbmap.Insert(&newserver)
				if derr != nil {
					fmt.Println("server insert error======>")
				}
    		}
		 	for x, p := range serverlist {
          		log.Printf("    %d: %v\n", x, p)
          		node_err := servers.InstallNode(nodeip, nodetype, p.Name)
				fmt.Printf("%s", node_err)
				if node_err != nil {
					result["success"] = false
				} else {
					uerr := updateServer(nodeip, p.Name)
				    if uerr != nil {
						fmt.Println("server insert error======>")
						result["success"] = false
					}
					result["success"] = true
					}
    		}
		 	var haserver orm.HAServers
		 	nerr := dbmap.SelectOne(&haserver, "select * from haservers where NodeIP2=?", nodeip)
			if nerr == nil {
				result["success"] = false
			} else {
				masterproxyerr := servers.InstallProxy(&haserver, "MASTER")
				if masterproxyerr != nil {
					result["success"] = false
				} 
				
				jsonMsg, _ := json.Marshal(haserver)
				
		 		url := "http://" + nodeip + ":8078/servernodes/haproxy/install"
		 		res, rerr := http.Post(url, "string", bytes.NewBufferString(string(jsonMsg)))		
	    		if rerr != nil {
		  	 	 	result["success"] = false
	    		} else {
		 			if res.StatusCode > 299 {
						result["success"] = false
					} else {
						result["success"] = true
					}
	  			}
	    	}
		 }	 
	}
}*/


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

func (this *ServerRouter) HAInstall() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	servername := this.Ctx.Input.Param(":name")
	
    //servername := "HAINSTALL"
	err := servers.InstallServers(servername)
	fmt.Printf("%s", err)
	if err != nil {
		result["success"] = false
	} else {
		result["success"] = true
	}
}

func (this *ServerRouter) ProxyInstall() {
	result := map[string]interface{}{
		"success": false,
	}

	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	req := this.Ctx.Request     //in beego this.Ctx.Request points to the Http#Request
	p := make([]byte, req.ContentLength)    
	ps, _ := this.Ctx.Request.Body.Read(p)
	fmt.Println(ps)
	
	var r orm.HAServers
    if haerr := json.Unmarshal(p, &r); haerr != nil {
    	result["success"] = false
    } else {
		slaveproxyerr := servers.InstallProxy(&r, "SLAVE")
		if slaveproxyerr != nil {
			result["success"] = false
		} else {
			result["success"] = true
		} 
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
	err := dbmap.SelectOne(&server, "select * from servers where Stype='MASTER' and Name=?", servername)

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

type Device struct {
	Name  		string  `json:"name"` 
	Size  		string  `json:"size"`
	State 		string  `json:"state"`
	MountPoint	string  `json:"mountpoint"`
}

type DeviceList struct {
	//Key		string  `json:"key"` 
	Key		*Device  `json:"key"` 
}

type Response struct {
	Data		[]*DeviceList  `json:"data"` 
	Host        string         `json:host"`
	Success		bool           `json:"success"` 
}

func (this *ServerRouter) GetHAOptions() {
	
	result := map[string]interface{}{
		"success": false,
	}
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	masterDevice := getDeviceDetails()
	
	ip := this.Ctx.Input.Param(":ip")
	url := "http://" + ip + ":8078/servernodes/devicedetails"
	res, err := http.Get(url)
	if err != nil {
		result["success"] = false
	} else {
		if res.StatusCode > 299 {
			result["success"] = false
		} else {
			haDevice, derr := ioutil.ReadAll(res.Body)
			if derr != nil {
		  		result["success"] = false
		  	} else {
		  		var r Response
    			if haerr := json.Unmarshal(haDevice, &r); haerr != nil {
        			result["success"] = false
    			} else {
    				fmt.Println(masterDevice)
    				disk1 := ""
    				disk2 := ""
    				for _, mk := range masterDevice {
    					for _, hav := range r.Data {
    						if mk.Key.Size == hav.Key.Size && mk.Key.State != "running" && hav.Key.State != "running" && mk.Key.MountPoint == "" && hav.Key.MountPoint == "" {
    						    if mk.Key.Size[len(mk.Key.Size)-1:] == "G" && hav.Key.Size[len(hav.Key.Size)-1:] == "G" {
    								disk1 = "/dev/"+mk.Key.Name
    								disk2 = "/dev/"+hav.Key.Name
    							}	
    						}
    				    }
    				}
    				details := getMasterDetails()
    				if details != "" {
    					detail := strings.Split(details, "=-=")
    					db := orm.OpenDB()
						dbmap := orm.GetDBMap(db)
    					newhaserver := orm.HAServers{NodeIP1: detail[0], NodeHost1: detail[1], NodeDisk1: disk1, NodeIP2: ip, NodeHost2: r.Host, NodeDisk2: disk2 }
						orm.ConnectToTable(dbmap, "haservers", newhaserver)
						derr := dbmap.Insert(&newhaserver)
						if derr != nil {
							fmt.Println("HA server insert error======>")
							result["success"] = false
						}
						result["success"] = true
    				} else {
    					result["success"] = false
    				}
    				
    			}		  		
		  	}
		}
	 }
	res.Body.Close()
}

func getMasterDetails() string {
	var e exec.OsExecutor
	var b bytes.Buffer
	var c bytes.Buffer
	var commandWords []string
	cmd := "bash conf/trusty/ip.sh"
	commandWords = strings.Fields(cmd)
   	err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
   	if err != nil {
    	return ""
    } else {
    	cmd1 := "hostname"
		commandWords = strings.Fields(cmd1)
   		err := e.Execute(commandWords[0], commandWords[1:], nil, &c, &c)
   		if err != nil {
    		return ""
	    } else {
    		return strings.TrimSpace(b.String())+"=-="+strings.TrimSpace(c.String())
    	}    		
    }
}

func (this *ServerRouter) HADeviceDetails() {
	var e exec.OsExecutor
	var b bytes.Buffer
	var commandWords []string
	result := map[string]interface{}{
		"success": false,
		"host": "",
	}
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	cmd := "hostname"
	commandWords = strings.Fields(cmd)
   	err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
   	if err != nil {
    	result["success"] = false 
    } else {
    	result["host"] = strings.TrimSpace(b.String())
    }
	
	haDevice := getDeviceDetails()
    result["data"] = haDevice
}



func getDeviceDetails() []*DeviceList {
	var e exec.OsExecutor
	var b bytes.Buffer
	var commandWords []string
	//devicelist := make([]string, 0)
	devicelist := make([]*DeviceList, 0)
	
	cmd := "lsblk -o NAME,SIZE,STATE,MOUNTPOINT -nl"
	commandWords = strings.Fields(cmd)
   	err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
   	if err != nil {
    	return devicelist
    } else {
    	/*testlines := "sda   232.9G running \n" +
                     "sda1   46.7G         / \n"+ 
					 "sda2    5.6G         [SWAP] \n" +
					 "sda3    244M         /boot \n" +
					 "sda4      1K         \n"+
					 "sda5   18.6G         /var \n"+
					 "sda6    4.7G         /usr \n"+
					 "sda7    4.7G         /home \n"+
					 "sda8    4.7G         /tmp \n"+
					 "sda9   18.6G         \n"+        
					 "sda10 129.2G         /storage1 \n"+
					 "sdb   232.9G running \n"+
					 "sdb1  116.4G         /storage2 \n"+
					 "sdb2  116.5G         /storage3 \n"+
					 "drbd0  18.6G         /var/lib/megam"      	
    	lines := strings.Split(testlines, "\n")*/
    	
    	lines := strings.Split(b.String(), "\n")
    	devicelist = make([]*DeviceList, len(lines))
    	for c, l := range lines {
    		line := strings.Split(l, " ")
    	    var s string
    		for _, n := range line {    		  
    		  if len(n) > 0 {
    		   	s = s + n + "-"
    		  }
    		}
    	  if len(s) > 0 { 	
    	  	listSplit := strings.Split(s, "-")
    	  	state := ""
    	  	mountpoint := "" 
    	  	if listSplit[2] == "running" {
    	  		state = listSplit[2]
    	  		mountpoint = listSplit[3]
    	  	} else {
    	  		state = ""
    	  		mountpoint = listSplit[2]
    	  	}
    	    devicelist[c] = &DeviceList{Key: &Device{Name: listSplit[0], Size: listSplit[1], State: state, MountPoint: mountpoint}}
    	  }
    	} 
    	fmt.Println(devicelist)
    }
	return devicelist
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
		t, err := tail.TailFile(filePath, tail.Config{Follow: true})
		fmt.Println(t)
		if err != nil {
			log.Printf("ERROR LOG READ ==> : %s", err.Error())
		}
		for line := range t.Lines {
			ip := strings.Split(line.Text, "=")
			err1 := dbmap.SelectOne(&node, "select * from nodes where IP=?", ip)
			if err1 != nil {
				result["ip"] = true
				result["ipvalue"] = line.Text
			}
			if len(line.Text) > 0 {
				t.Stop()
				err := os.Remove(filePath)
                
				if err != nil {
					fmt.Println(err)
				}
				return
			}
		}
}

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
		update_node := orm.Nodes{Id: node.Id, Install: true, IP: node.IP, HostName: node.HostName, InstallDate: node.InstallDate, UpdateDate: time.Format(layout)}
		orm.ConnectToTable(dbmap, "nodes", update_node)
		err2 := dbmap.Insert(&update_node)
	
		if err2 != nil {
			fmt.Println("node insert error======>")
			return err2
		}	
	return nil
}	

func updateServer(nodeip string, servername string) error {
        // insert rows - auto increment PKs will be set properly after the insert
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)

		server := orm.Servers{}
		err := dbmap.SelectOne(&server, "select * from servers where IP=? and Name=?", nodeip, servername)
		if err != nil {
			fmt.Println("select select error======>")
			return err
		}
		err3 := orm.DeleteRowFromServerNameAndIP(dbmap, servername, nodeip)
		if err3 != nil {
			log.Printf("server delete error")
			return err3
		}
		time := time.Now()
		update_server := orm.Servers{Id: server.Id, Name: server.Name, Install: true, IP: server.IP, Stype: server.Stype, HostName: server.HostName, InstallDate: server.InstallDate, UpdateDate: time.Format(layout)}
		orm.ConnectToTable(dbmap, "servers", update_server)
		err2 := dbmap.Insert(&update_server)
	
		if err2 != nil {
			fmt.Println("server insert error======>")
			return err2
		}	
	return nil
}	

