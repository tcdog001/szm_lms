package controllers

import (
	"LMS/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

type HeartData struct {
	Code int64 `json:"code"`
}

type HeartBeatController struct {
	beego.Controller
}

func (this *HeartBeatController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "login.html"
}

func (this *HeartBeatController) Post() {
	ret := HeartData{}
	//用户认证
	uname, pwd, ok := this.Ctx.Request.BasicAuth()
	if !ok {
		beego.Info("get client  Request.BasicAuth failed!")
		ret.Code = -1
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
		ret.Code = -1
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Info("user/pwd matched!")

	//接收请求信息
	beego.Debug("requestBody=", string(this.Ctx.Input.RequestBody))
	deviceinfo := models.Deviceinfo{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &deviceinfo)
	beego.Debug("deviceinfo=", deviceinfo)
	if err != nil {
		beego.Error(err)
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//心跳接收成功，更新设备状态
	deviceinfo.State = 1
	deviceinfo.LastKeepaliveTime = time.Now()
	if !models.UpdateDeviceStatus(&deviceinfo) {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//插入Listener中
	device := DeviceListener{
		State:         1,
		LastAliveTime: time.Now(),
	}
	Listener[deviceinfo.Mac] = device

	//返回给设备处理结果
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
}
