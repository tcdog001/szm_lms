package controllers

import (
	"LMS/models"
	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplNames = "login.html"
	this.Data["IsMatched"] = false
	ShowError := this.Input().Get("ShowError")
	this.Data["ShowError"] = ShowError
}

func (this *LoginController) Post() {
	//获取表单内容
	uname := this.Input().Get("uname")
	pwd := this.Input().Get("pwd")

	admin := models.Admininfo{
		Username: uname,
		Password: pwd,
	}
	//在数据库中做匹配
	ok := models.CheckAdmin(&admin)
	if ok {
		beego.Info("user/pwd matched!")

		this.SetSession("Admin", uname)
		if models.UpdateAdminStatus(&admin) {
			beego.Info("UpdateAdminStatus success!")
		} else {
			beego.Info("UpdateAdminStatus failed!")
		}
		beego.Info("Login success!")

		//登录成功，重定向到主页面
		this.Data["IsMatched"] = true
		path := "/home?CurUser=" + uname
		this.Redirect(path, 301)
		return
	}

	//登录失败，重定向到登录页面
	beego.Info("Login failed! Once again!")
	this.Data["IsMatched"] = false
	this.Redirect("/login/?ShowError=true", 302)
	return
}
