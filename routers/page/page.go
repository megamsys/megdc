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

package page

import (
	"strings"
	"encoding/json"
    "github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/routers/base"
    "bytes"
    "github.com/megamsys/libgo/exec"
   "net/http"
   "fmt"
   "io/ioutil"
)


// PageRouter serves home page.
type PageRouter struct {
	base.BaseRouter
}

// Get implemented dashboard page.
func (this *PageRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.Data["Username"] = this.GetUser()
	this.TplNames = "page/index.html" 
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
	
}


func (this *PageRouter) Dash() {
	this.Data["Username"] = this.GetUser()
	this.TplNames = "page/dash.html"
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
}

func (this *PageRouter) Master() {
	this.Data["Username"] = this.GetUser()
	this.TplNames = "page/master.html"
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
}

func (this *PageRouter) MasterDashboard() {
	var e exec.OsExecutor
	var serverlist []orm.Servers
	var servname string
	var b bytes.Buffer

	var commandWords []string
	result := map[string]interface{}{
		"success": false,
	}
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	_, err := dbmap.Select(&serverlist, "select * from servers where Stype='MASTER'") 
	
	if len(serverlist) == 0 || err != nil {
		result["success"] = false
	} else {
		for _, p := range serverlist {
       	   servname = servname + " " + p.Name
       	 }  
		//cmd := "/home/rajthilak/.rvm/rubies/ruby-2.2.0/bin/ruby conf/trusty/cib.rb" + strings.ToLower(servname)
		cmd := "ruby conf/trusty/cib.rb" + strings.ToLower(servname)
		commandWords = strings.Fields(cmd)
    	err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
       	if err != nil {
    		result["success"] = false
    	} else {
    		result["success"] = true
    	    result["data"] = b.String()
    	 }   
    }
}

func (this *PageRouter) HADashboard() {
	var serverlist []orm.Servers
	var servname string
	result := map[string]interface{}{
		"success": false,
	}
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	_, err := dbmap.Select(&serverlist, "select * from servers where Stype='HA'") 
	
	if len(serverlist) == 0 || err != nil {
		result["success"] = false
	} else {
		for _, p := range serverlist {
       	   servname = servname + " " + p.Name
       	 }  
		url := "http://" + serverlist[0].IP + ":8078/dashboard/ha/request"
	    res, rerr := http.Post(url, "string", bytes.NewBufferString(servname))		
		if rerr != nil {
			result["success"] = false
		} else {
			if res.StatusCode == 200 {
				resBody, derr := ioutil.ReadAll(res.Body)
		  		if derr != nil {
		  			result["success"] = false
		  		} else {		  		  
		  		    var dat map[string]interface{}
    				if jerr := json.Unmarshal(resBody, &dat); jerr != nil {
        				result["success"] = false
    				} else {
    					if dat["success"] == false {
		  					result["success"] = false
		  				} else {
    						result["success"] = true
							//c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
                			result["data"] = dat["data"]
                		}			
		  			} 		
    			}			
			} else {
				result["success"] = false
			}
		}
	   	res.Body.Close()
	 }
}

func (this *PageRouter) HADashboardRequest() {
	var e exec.OsExecutor
	var b bytes.Buffer
   // var servernames []orm.Servers
    var servname string
	var commandWords []string
	result := map[string]interface{}{
		"success": false,
	}
	req := this.Ctx.Request     //in beego this.Ctx.Request points to the Http#Request
	p := make([]byte, req.ContentLength)    
	ps, _ := this.Ctx.Request.Body.Read(p)
	fmt.Println(ps)
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()

    servname = strings.ToLower(string(p))
	//cmd := "/home/rajthilak/.rvm/rubies/ruby-2.2.0/bin/ruby conf/trusty/cib.rb" + servname
	cmd := "ruby conf/trusty/cib.rb" + servname
	commandWords = strings.Fields(cmd)
   	err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
   	if err != nil {
    	result["success"] = false
    } else {
    	result["success"] = true
  		//  c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
    	result["data"] = b.String()
    }
}

func (this *PageRouter) CSDashboard() {
	var nodelist []orm.Nodes
	
	result := map[string]interface{}{
		"success": false,
	}
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	_, err := dbmap.Select(&nodelist, "select * from nodes") 
	
	if len(nodelist) == 0 || err != nil {
		result["success"] = false
	} else {
		for _, n := range nodelist {
			url := "http://" + n.IP + ":8078/dashboard/cs/request"
	    	res, rerr := http.Get(url)			
			if rerr != nil {
				result["success"] = false
			} else {
				if res.StatusCode == 200 {
					resBody, derr := ioutil.ReadAll(res.Body)
		  			if derr != nil {
		  				result["success"] = false
		  			} else {		  		  
		  		 	   var dat map[string]interface{}
    					if jerr := json.Unmarshal(resBody, &dat); jerr != nil {
        					result["success"] = false
    					} else {
    						if dat["success"] == false {
		  						result["success"] = false
		  					} else {
    							result["success"] = true
								//c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
                				result["data"] = dat["data"]
                			}			
		  				} 		
    				}			
				} else {
					result["success"] = false
				}
			}
	   	res.Body.Close()
	  } 	
	 }
}

func (this *PageRouter) CSDashboardRequest() {
	var e exec.OsExecutor
	var b bytes.Buffer

	var commandWords []string
	result := map[string]interface{}{
		"success": false,
	}
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	//cmd := "/home/rajthilak/.rvm/rubies/ruby-2.2.0/bin/ruby conf/trusty/cib.rb opennebula_host ceph"
	cmd := "ruby conf/trusty/cib.rb opennebula_host ceph"
	commandWords = strings.Fields(cmd)
    err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
    if err != nil {
    	result["success"] = false
    } else {
    	fmt.Println("-------------------------------------------")
    	fmt.Println(b.String())
    	fmt.Println(err)
    	result["success"] = true
  		//  c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
    	result["data"] = b.String()
    }
}



