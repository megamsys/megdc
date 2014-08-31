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

package auth

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/megamsys/cloudinabox/modules/auth"
	"github.com/megamsys/cloudinabox/modules/utils"
	"github.com/megamsys/cloudinabox/routers/base"
	"net/http"
	"strings"
)

// LoginRouter serves login page.
type LoginRouter struct {
	base.BaseRouter
}

// Get implemented login page.
func (this *LoginRouter) Get() {
	loginRedirect := strings.TrimSpace(this.GetString("to"))
	if len(this.Ctx.GetCookie("remember")) > 0 {
		this.Data["Username"] = this.Ctx.GetCookie("user_name")
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
	response, _ := this.Auth(client, data)
    fmt.Println(response)
    if response != nil {
    	if response.StatusCode > 399 && response.StatusCode < 498 {
		   this.FlashWrite("LoginError", "true")
		   this.Redirect("/", 302)
	   } else if response.StatusCode > 499 {
		   this.FlashWrite("ServerError", "true")
		   this.Redirect("/", 302)
	   } else {
		this.LoginUser(data, true)
		this.Redirect("/index", 302)
	  }
    } else {
    	this.FlashWrite("ServerError", "true")
		this.Redirect("/", 302)
    }

//	if err != nil {
//		this.FlashWrite("LoginError", "true")
//		this.Redirect("/", 302)
//	} else {
//		this.LoginUser(data, true)
//		this.Redirect("/index", 302)
//	}

}

// Logout implemented user logout page.
func (this *LoginRouter) Logout() {
	auth.LogoutUser(this.Ctx)

	// write flash message
	this.FlashWrite("HasLogout", "true")

	this.Redirect("/", 302)
}
