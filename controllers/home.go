package controllers

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	beego.Debug("session id=", session)
	CurUser := this.Input().Get("CurUser")
	beego.Debug("CurUser=", CurUser)

	this.TplNames = "home.html"
	this.Data["CurUser"] = CurUser
}

func (this *HomeController) Post() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "login.html"
}
