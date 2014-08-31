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

// Package routers implemented controller methods of beego.
package base

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"github.com/megamsys/cloudinabox/models/orm"
	"github.com/megamsys/cloudinabox/modules/auth"
	"github.com/megamsys/cloudinabox/modules/utils"
//	"html/template"
	"net/url"
	"strings"
	"time"
	//	"strconv"
	"bytes"
	"net/http"
)

var (
	AppVer string
	IsPro  bool
)

var langTypes []*langType // Languages are supported.

// langType represents a language type.
type langType struct {
	Lang, Name string
}

// baseRouter implemented global settings for all other routers.
type BaseRouter struct {
	beego.Controller
	i18n.Locale
	//	User orm.Users
	IsLogin bool
}

// Prepare implemented Prepare method for baseRouter.
func (this *BaseRouter) Prepare() {
	// Setting properties.
	this.Data["AppVer"] = AppVer
	this.Data["IsPro"] = IsPro

	this.Data["PageStartTime"] = time.Now()

	//this.StartSession()
	// check flash redirect, if match url then end, else for redirect return
	if match, redir := this.CheckFlashRedirect(this.Ctx.Request.RequestURI); redir {
		return
	} else if match {
		this.EndFlashRedirect()
	}

	// pass xsrf helper to template context
	//xsrfToken := this.XsrfToken()
	//this.Data["xsrf_token"] = xsrfToken
	//this.Data["xsrf_html"] = template.HTML(this.XsrfFormHtml())

	// read flash message
	beego.ReadFromRequest(&this.Controller)

}

// setLangVer sets site language version.
func (this *BaseRouter) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	this.Lang = lang
	this.Data["Lang"] = curLang.Lang
	this.Data["CurLang"] = curLang.Name
	this.Data["RestLangs"] = restLangs

	return isNeedRedir
}

func (this *BaseRouter) Auth(client *utils.Client, data *utils.User) (*http.Response, error) {

	//we need to move into a struct
	tmpinp := map[string]string{
		"email":     data.Username,
		"api_key":   data.Api_key,
		"authority": "user",
	}

	//and this as well.
	jsonMsg, err := json.Marshal(tmpinp)

	if err != nil {
		return nil, err
	}

	authly, err := utils.NewAuthly("/auth", jsonMsg, data)
	if err != nil {
		return nil, err
	}

	url, err := utils.GetURL("/auth")
	if err != nil {
		return nil, err
	}

	fmt.Println("==> " + url)
	authly.JSONBody = jsonMsg

	err = authly.AuthHeader()
	if err != nil {
		return nil, err
	}
	client.Authly = authly

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonMsg))
	if err != nil {
		return nil, err
	}

	//resp, err := client.Do(request)
	//if err != nil {
	////	return err
	//}
	// fmt.Println(strconv.Itoa(resp.StatusCode) + " ....code")
	//if resp.StatusCode == http.StatusNoContent {
	//    fmt.Println("Service successfully updated.")
	//fmt.Fprintln(ctx.Stdout, "Service successfully updated.")
	//}
	return client.Do(request)
}

// check if not login then redirect
func (this *BaseRouter) CheckLoginRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			// custom redirect url
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !this.IsLogin {
		if len(redirect_to) == 0 {
			req := this.Ctx.Request
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/?to=" + url.QueryEscape(redirect_to)
		this.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && this.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		this.Redirect(redirect_to, code)
		return true
	}
	return false
}

// read beego flash message
func (this *BaseRouter) FlashRead(key string) (string, bool) {
	if data, ok := this.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash error message
func (this *BaseRouter) FlashErrorWrite(key string, value string) {
	flash := beego.NewFlash()
	flash.Error(value)
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

// write beego flash message
func (this *BaseRouter) FlashWrite(key string, value string) {
	flash := beego.NewFlash()
	//	flash.Notice(value)
	flash.Data[key] = value
	flash.Store(&this.Controller)
}

func (this *BaseRouter) LoginUser(user *utils.User, remember bool) string {
	loginRedirect := strings.TrimSpace(this.Ctx.GetCookie("login_to"))
	if utils.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	} else {
		this.Ctx.SetCookie("login_to", "", -1, "/")
	}
	//set cookie remember true
	this.Ctx.SetCookie("remember", "true")

	// login user
	auth.LoginUser(user, this.Ctx, remember)

	// write user details in database
	// insert rows - auto increment PKs will be set properly after the insert
	db := orm.OpenDB()
	dbmap := orm.GetDBMap(db)
	newuser := orm.NewUser(user)
	orm.ConnectToTable(dbmap, "users", newuser)
	err := dbmap.Insert(&newuser)
	if err != nil {
		fmt.Println("insert error======>")
		return ""
	}
	defer db.Close()
	//	this.setLangCookie(i18n.GetLangByIndex(user.Lang))

	return loginRedirect
}

func (this *BaseRouter) GetUser() string {
	return this.Ctx.GetCookie("user_name")
}

// check flash redirect, ensure browser redirect to uri and display flash message.
func (this *BaseRouter) CheckFlashRedirect(value string) (match bool, redirect bool) {
	v := this.GetSession("on_redirect")
	if params, ok := v.([]interface{}); ok {
		if len(params) != 5 {
			this.EndFlashRedirect()
			goto end
		}
		uri := utils.ToStr(params[0])
		code := 302
		if c, ok := params[1].(int); ok {
			if c/100 == 3 {
				code = c
			}
		}
		flag := utils.ToStr(params[2])
		flagVal := utils.ToStr(params[3])
		times := 0
		if v, ok := params[4].(int); ok {
			times = v
		}

		times += 1
		if times > 3 {
			// if max retry times reached then end
			this.EndFlashRedirect()
			goto end
		}

		// match uri or flash flag
		if uri == value || flag == value {
			match = true
		} else {
			// if no match then continue redirect
			this.FlashRedirect(uri, code, flag, flagVal, times)
			redirect = true
		}
	}
end:
	return match, redirect
}

// set flash redirect
func (this *BaseRouter) FlashRedirect(uri string, code int, flag string, args ...interface{}) {
	flagVal := "true"
	times := 0
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			flagVal = v
		case int:
			times = v
		}
	}

	if len(uri) == 0 || uri[0] != '/' {
		panic("flash redirect only support same host redirect")
	}

	params := []interface{}{uri, code, flag, flagVal, times}
	this.SetSession("on_redirect", params)

	this.FlashWrite(flag, flagVal)
	this.Redirect(uri, code)
}

// clear flash redirect
func (this *BaseRouter) EndFlashRedirect() {
	this.DelSession("on_redirect")
}


