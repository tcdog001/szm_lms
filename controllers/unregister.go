package controllers

import (
	"LMS/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

type UnRegisterData struct {
	Code int64 `json:"code"`
}

type UnRegisterController struct {
	beego.Controller
}

func (this *UnRegisterController) Get() {
	this.TplNames = "login.html"
}

/*返回值
格式：{"code":0/-1/-2}
( 0) success
(-1) user/password error
(-2) other error
*/
func (this *UnRegisterController) Post() {
	ret := UnRegisterData{}
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

	//获取请求的json数据
	beego.Info("request body=", string(this.Ctx.Input.RequestBody))
	deviceinfo := models.Deviceinfo{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &deviceinfo)
	if err != nil {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//更新数据库中设备的状态为离线,并删除状态监听器中对应项
	deviceinfo.State = 0
	models.UpdateDeviceStatus(&deviceinfo)
	delete(Listener, deviceinfo.Mac)

	//返回注销成功
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	beego.Info(string(writeContent))
	return
}
