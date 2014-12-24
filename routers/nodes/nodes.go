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

package nodes

import (
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/routers/base"
//	"fmt"
	"encoding/json"
)

var serversList = [...]string{"MEGAM", "COBBLER", "OPENNEBULA", "OPENNEBULAHOST"}


// PageRouter serves home page.
type NodesRouter struct {
	base.BaseRouter
}

// Get implemented dashboard page.
func (this *NodesRouter) Nodes() {
	this.Data["Username"] = this.GetUser()
	this.TplNames = "page/nodes.html"
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
}

func (this *NodesRouter) Ha() {
	this.Data["Username"] = this.GetUser()
	this.TplNames = "page/ha.html"
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
}


// Get implemented dashboard page.
func (this *NodesRouter) GetNodes() {
	var nodes []*orm.Nodes
	var nodesempty []*orm.Nodes
//	var nodes_output []string
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
	_, err := dbmap.Select(&nodes, "select * from nodes")
	if err != nil {
		result["success"] = false
	    result["data"] = nodesempty
	} else {
	   jsonMsg, _ := json.Marshal(nodes)
	   result["success"] = true
	   result["data"] = string(jsonMsg)
  }
}




