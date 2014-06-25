package controllers

import (
	"github.com/astaxie/beego"  
)

type LoginController struct {
    beego.Controller
}

func (this *LoginController) Get() {
    this.TplNames = "login.html"
    //form := auth.LoginForm{}
   // this.Ctx.WriteString("hello world")
}

type SignInController struct {
  username string
  api_key string
  beego.Controller
  }
  
func (this *SignInController) Post() {
 this.username = this.Ctx.Input.Param(":email")
 this.api_key = this.Ctx.Input.Param(":pass")  
 this.Ctx.WriteString(this.username)
 }