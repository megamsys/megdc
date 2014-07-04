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

package auth

import (
	"strings"
	"github.com/megamsys/cloudinabox/modules/utils"
	"github.com/megamsys/cloudinabox/modules/auth"
    "github.com/megamsys/cloudinabox/routers/base"
    _ "github.com/mattn/go-sqlite3"
    "net/http"
    _ "github.com/lib/pq"
    "fmt"
)

// LoginRouter serves login page.
type LoginRouter struct {
	base.BaseRouter
}

// Get implemented login page.
func (this *LoginRouter) Get() {   
	loginRedirect := strings.TrimSpace(this.GetString("to"))
    if len(this.Ctx.GetCookie("remember")) > 0 {
    	this.Data["Username"] = "Megam"
    	this.Redirect("/index", 302)
    } else {
    	this.Data["IsLoginPage"] = true
    	this.TplNames = "auth/login.html" 
    	this.Ctx.SetCookie("login_to", loginRedirect, 0, "/")
    }	
}

// Login implemented user login.
func (this *LoginRouter) Login() {
	this.Data["IsLoginPage"] = true
	this.TplNames = "auth/login.html"
	if this.CheckLoginRedirect(false) {
		return
	}
	
    data := &utils.User{this.GetString("username"), this.GetString("password")}
    client := utils.NewClient(&http.Client{}, data)
	_, err := this.Auth(client, data)

	if err != nil {
		this.FlashWrite("Email and password was wrong. Please re-entry the fields.", "true")
		this.Redirect("/", 302)
		fmt.Println(err)
	} else {
		this.LoginUser(data, true)
		this.Redirect("/index", 302)
	}
		
}

// Logout implemented user logout page.
func (this *LoginRouter) Logout() {
	auth.LogoutUser(this.Ctx)

	// write flash message
	this.FlashWrite("HasLogout", "true")

	this.Redirect("/", 302)
}
