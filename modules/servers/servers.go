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
	"errors"
	"time"
	"reflect"
    "unsafe"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/megamsys/cloudinabox/app"
	"github.com/megamsys/cloudinabox/models/orm"
	"net/http"
)

const layout = "Jan 2, 2006 at 3:04pm (MST)"

func InstallServers(serverName string) error {
	var err error
	switch serverName {
	case "MEGAM":
		err = app.MegamInstall()
		if err != nil {
			fmt.Printf("Error: Install error for [%s]", serverName)
			fmt.Println(err)
			return err
		}
	case "COBBLER":
		err = app.CobblerInstall()
		if err != nil {
			fmt.Printf("Error: Install error for [%s]", serverName)
			fmt.Println(err)
			return err
		}
	case "OPENNEBULA":
		err = app.NebulaInstall()
		if err != nil {
			fmt.Printf("Error: Install error for [%s]", serverName)
			fmt.Println(err)
			return err
		}
	case "OPENNEBULAHOST":
		err = app.OpenNebulaHostMasterInstall()
		if err != nil {
			fmt.Printf("Error: Install error for [%s]", serverName)
			fmt.Println(err)
			return err
		}
	case "NODEINSTALL":
		err = app.OpenNebulaHostNodeInstall()
		if err != nil {
			fmt.Printf("Error: Install error for [%s]", serverName)
			fmt.Println(err)
			return err
		}	
     case "HAINSTALL":
		err = app.HANodeInstall()
		if err != nil {
			fmt.Printf("Error: Install error for [%s]", serverName)
			fmt.Println(err)
			return err
		}	
	}
	return nil
}

func CheckHAInstall(servername string) error {
	var serverlist []orm.Servers
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	_, err := dbmap.Select(&serverlist, "select * from servers where Stype='HA' GROUP BY IP") 
	if err != nil {
		return err
	} else {
		if len(serverlist) > 0 {
			for _, v := range serverlist {
				newserver := orm.NewServer(servername, v.IP, "HA", "")
				orm.ConnectToTable(dbmap, "servers", newserver)
				derr := dbmap.Insert(&newserver)
				if derr != nil {
					fmt.Println("server insert error======>")
				}
				
					url := "http://" + v.IP + ":8078/servernodes/ha/" + servername + "/install"
	    			res, err := http.Get(url)
	    			if err != nil {
		    			return err
	   				 } else {
		 				if res.StatusCode > 299 {
							return errors.New(res.Status)
						}
		 				uerr := updateServer(v.IP, servername)
						if uerr != nil {
							fmt.Println("server insert error======>")
							return uerr
						} 
	  				}
			}
		}
  	}
	return nil
}

type Response struct {
	Success		bool           `json:"success"` 
	Error       string         `json:"errordata"`
}

func BytesToString(b []byte) string {
    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
    sh := reflect.StringHeader{bh.Data, bh.Len}
    return *(*string)(unsafe.Pointer(&sh))
}

func InstallNode(nodeip string, nodetype string, name string) error {
	if nodetype == "COMPUTE" {
	url := "http://" + nodeip + ":8078/servernodes/nodes/install"
	res, err := http.Get(url)
    	
	if err != nil {
		return err
	} else {
		if res.StatusCode > 299 {
			return errors.New(res.Status)
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			var r Response
    		if haerr := json.Unmarshal(body, &r); haerr != nil {
    			return haerr
    		} else {
    			if r.Success {
    				err = app.SCPSSHInstall()
					return err
    			} else {
    				return errors.New(r.Error)
    			}
    		}			
		}
	 }
	} else {
		url := "http://" + nodeip + ":8078/servernodes/ha/" + name + "/install"
	    res, err := http.Get(url)
	    if err != nil {
		    return err
	    } else {
		 if res.StatusCode > 299 {
			return errors.New(res.Status)
		} else {
			return nil
		}
	  }
	}
}

func InstallProxy(haserver *orm.HAServers, Stype string) error {
	cib := &app.CIB{}
	if Stype == "MASTER" {
		cib = &app.CIB{LocalIP: haserver.NodeIP1, LocalHost: haserver.NodeHost1, LocalDisk: haserver.NodeDisk1, RemoteIP: haserver.NodeIP2, RemoteHost: haserver.NodeHost2, RemoteDisk: haserver.NodeDisk2, Master: true}
	} else {
		cib = &app.CIB{LocalIP: haserver.NodeIP2, LocalHost: haserver.NodeHost2, LocalDisk: haserver.NodeDisk2, RemoteIP: haserver.NodeIP1, RemoteHost: haserver.NodeHost1, RemoteDisk: haserver.NodeDisk1, Master: false}
	}
	err := app.HAProxyInstall(cib, Stype)
		if err != nil {
			fmt.Printf("Error: Install error for HAProxy")
			fmt.Println(err)
			return err
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
			fmt.Println("server delete error")
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
