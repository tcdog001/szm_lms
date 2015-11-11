package controllers

import (
	"LMS/models"
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

var base *base64.Encoding

func init() {
	//base64编码
	base = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
}

type CommandData struct {
	Id  int64  `json:"id"`
	Cmd string `json:"cmd"`
}

type GetCommandController struct {
	beego.Controller
}

func (this *GetCommandController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "login.html"
}

func (this *GetCommandController) Post() {
	ret := CommandData{}
	//身份认证
	uname, pwd, ok := this.Ctx.Request.BasicAuth()
	if !ok {
		beego.Info("get client  Request.BasicAuth failed!")
		ret.Cmd = ""
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	user := models.Userinfo{
		Username: uname,
		Password: pwd,
	}
	ok = models.CheckAccount(&user)
	if !ok {
		beego.Info("user/pwd not matched!")
		ret.Cmd = ""
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Info("user/pwd matched!")

	//获取请求信息
	deviceinfo := models.Deviceinfo{
		State:             1,
		LastKeepaliveTime: time.Now(),
	}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &deviceinfo)
	if err != nil {
		ret.Cmd = ""
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//设备取命令的动作做为一次心跳，更新设备状态
	models.UpdateDeviceStatus(&deviceinfo)

	//插入Listener中
	device := DeviceListener{
		State:         1,
		LastAliveTime: time.Now(),
	}
	Listener[deviceinfo.Mac] = device

	//取命令
	command := models.Command{
		Executed: false,
	}
	json.Unmarshal(this.Ctx.Input.RequestBody, &command)
	ok, com := models.GetCommand(&command)
	if !ok {
		//该设备没有对应的命令，返回命令为空
		ret.Cmd = ""
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Info(com)

	//设备取出命令，更新状态为（等待下发命令）
	com.Executed = true
	ok = models.UpdateDeviceCommand(com)
	if !ok {
		ret.Cmd = ""
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//更新操作记录状态为(已执行)
	/*
		record := models.OperationRecord{
			Mac:      com.Mac,
			Command:  com.Command,
			Executed: true,
			ExecTime: time.Now(),
		}*/
	record := com.Operecord
	record.Executed = true
	record.ExecTime = time.Now()
	ok = models.UpdateOperationRecord(record)
	if !ok {
		ret.Cmd = ""
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//将命令发送给设备
	ret.Id = com.Id
	ret.Cmd = base.EncodeToString([]byte(com.Command))
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
}
