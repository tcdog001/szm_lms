package controllers

import (
	"LMS/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

type GetScriptData struct {
	Code int64  `json:"code"`
	Path string `json:"path"`
}

type GetScriptController struct {
	beego.Controller
}

var URL = "http://lms1.autelan.com:8080/" + ScriptFolder

func (this *GetScriptController) Get() {
	this.TplNames = "login.html"
}

/*返回值
格式：{"code":0/-1/-2}
( 0) success
(-1) user/password error
(-2) other error
*/

func (this *GetScriptController) Post() {
	ret := GetScriptData{}
	//用户身份认证
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

	//获取设备post数据
	beego.Info("request body=", string(this.Ctx.Input.RequestBody))
	script := models.Script{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &script)
	if err != nil {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//获取对应脚本信息
	sc, ok := models.GetScript(&script)
	if !ok {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Debug(sc)
	//修改状态
	sc.Downloaded = true
	sc.DowonlodTime = time.Now()
	ok = models.UpdateScript(sc)
	if !ok {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	} else {
		ret.Path = URL + sc.FilePath
	}
	//更新操作记录状态为(已执行)
	record := models.OperationRecord{
		Mac:      sc.Mac,
		Command:  "",
		Script:   sc.FilePath,
		Executed: true,
		ExecTime: time.Now(),
	}
	ok = models.UpdateOperationRecord(&record)
	if !ok {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	//返回成功
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	beego.Info(string(writeContent))
	return
}
