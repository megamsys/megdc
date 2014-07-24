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
//	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/context"
//	"strings"
	// "time"

	"github.com/astaxie/beego"
//	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"

	"github.com/megamsys/cloudinabox/modules/utils"
	"github.com/megamsys/cloudinabox/setting"
)

var (
	user_name string
)

func GetUserIdFromSession(sess session.SessionStore) string {
	if user_name, ok := sess.Get("auth_user_name").(string); ok && len(user_name) > 0 {
		fmt.Println("auth user entry==========")
		return user_name
	}
	return ""
}

// get user if key exist in session
func GetUserFromSession(sess session.SessionStore) bool {
	name := GetUserIdFromSession(sess)
	if name != "" {
	//	u := models.User{Id: id}
	//	if u.Read() == nil {
	//		*user = u
			return true
	//	}
	}

	return false
}

func DeleteRememberCookie(ctx *context.Context) {
	ctx.SetCookie(setting.CookieUserName, "", -1)
	ctx.SetCookie(setting.CookieRememberName, "", -1)
	ctx.SetCookie("remember", "", -1)
}

func LoginUserFromRememberCookie(user *utils.User, ctx *context.Context) (success bool) {
	userName := ctx.GetCookie(setting.CookieUserName)
	if len(userName) == 0 {
		return false
	}

	defer func() {
		if !success {
			DeleteRememberCookie(ctx)
		}
	}()
	

	//secret := utils.EncodeMd5(user.Rands + user.Password)
	secret := utils.EncodeMd5("abcd" + "megam")
	value, _ := ctx.GetSecureCookie(secret, setting.CookieRememberName)
	if value != userName {
		return false
	}

	LoginUser(user, ctx, true)

	return true
}

// login user
func LoginUser(user *utils.User, ctx *context.Context, remember bool) {
	// werid way of beego session regenerate id...
	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
	ctx.Input.CruSession = beego.GlobalSessions.SessionRegenerateId(ctx.ResponseWriter, ctx.Request)
	ctx.Input.CruSession.Set("auth_user_name", user.Username)
	days := 86400 * setting.LoginRememberDays
    ctx.SetCookie("user_name", user.Username, days)
	if remember {
		WriteRememberCookie(user, ctx)
	}
}

func WriteRememberCookie(user *utils.User, ctx *context.Context) {
	secret := utils.EncodeMd5(user.Username + user.Api_key)
	days := 86400 * setting.LoginRememberDays
	ctx.SetCookie(setting.CookieUserName, user.Username, days)
	ctx.SetSecureCookie(secret, setting.CookieRememberName, user.Username, days)
}

// logout user
func LogoutUser(ctx *context.Context) {
	DeleteRememberCookie(ctx)
	ctx.Input.CruSession.Delete("auth_user_name")
	ctx.Input.CruSession.Flush()
	beego.GlobalSessions.SessionDestroy(ctx.ResponseWriter, ctx.Request)
}


