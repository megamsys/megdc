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

package page

import (
//	"strings"
//	"github.com/megamsys/cloudinabox/modules/utils"
    "github.com/megamsys/cloudinabox/routers/base"
 //   "net/http"
 //   "regexp"
 //  "fmt"
)

// PageRouter serves home page.
type PageRouter struct {
	base.BaseRouter
}

// Get implemented dashboard page.
func (this *PageRouter) Get() {
	this.Data["IsLoginPage"] = true
	this.Data["Username"] = "Megam"
	this.TplNames = "page/dashboard.html" 
	if len(this.Ctx.GetCookie("remember")) == 0 {
		this.Redirect("/", 302)
	}
}