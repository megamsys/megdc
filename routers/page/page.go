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
    "github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/routers/base"
 //   "os/exec"
    "bytes"
    "github.com/megamsys/libgo/exec"
   "net/http"
 //   "regexp"
   "fmt"
   "io/ioutil"
)


// PageRouter serves home page.
type PageRouter struct {
	base.BaseRouter
}

// Get implemented dashboard page.
func (this *PageRouter) Get() {
	//result := map[string]interface{}{
	//	"success": true,
	//}

	//defer func() {
	///	this.Data["json"] = result
	//	this.ServeJson()
	//}()
	
	
	//servers := new(orm.Servers)
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
	var b bytes.Buffer

	var commandWords []string
	result := map[string]interface{}{
		"success": false,
	}
	
	defer func() {
		this.Data["json"] = result
		this.ServeJson()
	}()
	
	cmd := "/home/rajthilak/.rvm/rubies/ruby-2.2.0/bin/ruby conf/trusty/cib.rb megam cobbler"
	commandWords = strings.Fields(cmd)
    err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
    
    fmt.Println("-------------------------------------------")
    fmt.Println(b.String())
    fmt.Println(err)
    c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
    result["data"] = c
}

func (this *PageRouter) HADashboard() {
	var serverlist []orm.Servers
	
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
		url := "http://" + serverlist[0].IP + ":8078/dashboard/ha/request"
	    	res, err := http.Get(url)
				robots, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
					if err != nil {
					fmt.Println(err)
					}
				//	fmt.Printf("%s", string(robots).data)
	   	 if err != nil {
			result["success"] = false
	   	} else {
		 	if res.StatusCode == 200 {
				result["success"] = true
				c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
                result["data"] = c
			} else {
				result["success"] = false
			}
	  	}
	 }
}

func (this *PageRouter) HADashboardRequest() {
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
	
	cmd := "/home/rajthilak/.rvm/rubies/ruby-2.2.0/bin/ruby conf/trusty/cib.rb megam cobbler"
	commandWords = strings.Fields(cmd)
    err := e.Execute(commandWords[0], commandWords[1:], nil, &b, &b)
    
    fmt.Println("-------------------------------------------")
    fmt.Println(b.String())
    fmt.Println(err)
    result["success"] = true
    c := "{\"packages\":{\"megam\":{\"megamcommon\":\"true\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"megamsnowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"true\",\"sqlite3\":\"true\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"true\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\",\"debmirror\":\"true\"}},\"services\":{\"megam\":{\"megamcommon\":\"false\",\"megamcib\":\"false\",\"megamcibn\":\"false\",\"megamnilavu\":\"false\",\"snowflake\":\"true\",\"megamgateway\":\"false\",\"megamd\":\"false\",\"megamchefnative\":\"false\",\"megamanalytics\":\"false\",\"megamdesigner\":\"false\",\"megammonitor\":\"false\",\"riak\":\"true\",\"rabbitmq-server\":\"true\",\"nodejs\":\"false\",\"sqlite3\":\"false\",\"ruby2.0\":\"false\",\"openjdk-7-jdk\":\"false\"},\"cobbler\":{\"cobbler\":\"true\",\"dnsmasq\":\"true\",\"apache2\":\"true\"}}}"
    result["data"] = c
   /*return &http.Response{
    Status:     "200 OK",
    StatusCode: 200,
    Body: ioutil.NopCloser(bytes.NewBufferString(c)),
   }*/
}




