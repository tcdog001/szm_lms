package controllers

import (
	"LMS/models"
	"github.com/astaxie/beego"
	//"strconv"
	//"time"
)

type CommandController struct {
	beego.Controller
}

func (this *CommandController) Get() {
	//检查登录状态
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "command.html"

	//获取传递数据
	mac := this.Input().Get("mac")
	this.Data["Mac"] = mac
	CurUser := this.Input().Get("CurUser")
	this.Data["CurUser"] = CurUser
}

func (this *CommandController) Post() {
	//检查登录状态
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	//获取表单信息
	admin := this.Input().Get("CurUser")
	beego.Debug("admin=", admin)
	mac := this.Input().Get("mac")
	beego.Debug("mac=", mac)
	commandContent := this.Input().Get("commandContent")
	beego.Debug("commandContent=", commandContent)

	/*
		//判断该设备命令状态,如果存在命令未执行,则拒绝添加
		exist := models.CheckCommandExist(&command)
		if exist {
			beego.Debug("the device command exist!")
			path := "/deviceinfo?CurUser=" + admin + "&CommandExist=" + strconv.FormatBool(true)
			this.Redirect(path, 302)
			return
		}
		beego.Debug("the device command not exist! can add command")
	*/

	//数据库中添加管理员操作记录
	record := models.OperationRecord{
		Operator: admin,
		Mac:      mac,
		Command:  commandContent,
		Executed: false,
	}
	ok := models.AddOperationRecord(&record)
	if ok {
		//下发命令，设置状态为（待取）
		command := models.Command{
			Mac:       mac,
			Command:   commandContent,
			Executed:  false,
			Operecord: &record,
		}
		models.AddDeviceCommand(&command)
	}

	//返回设备页面
	path := "/deviceinfo?CurUser=" + admin
	this.Redirect(path, 302)
	return
}
