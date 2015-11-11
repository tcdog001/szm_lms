package controllers

import (
	"LMS/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

var alarmsCount int64

type AlarmData struct {
	Code int64 `json:"code"`
}

type AlarmController struct {
	beego.Controller
}

func (this *AlarmController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "alarm.html"
	admin := this.Input().Get("CurUser")
	beego.Debug("admin=", admin)
	this.Data["CurUser"] = admin

	mac := this.Input().Get("Mac")
	beego.Debug("mac=", mac)

	//重置参数
	alarmsCount = 0
	totalPages = 0
	curPage = 1
	listCount = 10
	curCount = 0

	//获取每页显示条数
	listcount := this.Input().Get("ListCount")
	if listcount != "" {
		listCount, _ = strconv.ParseInt(listcount, 10, 0)
		beego.Debug("listCount=", listCount)
	}

	//获取记录总条数
	alarmscount := this.Input().Get("AlarmsCount")
	if alarmscount != "" {
		alarmsCount, _ = strconv.ParseInt(alarmscount, 10, 0)
		beego.Debug("recordsCount=", alarmsCount)
	} else {
		alarmsCount = models.GetAlarmCount()
	}
	//计算总页数
	if alarmsCount%listCount > 0 {
		totalPages = alarmsCount/listCount + 1
	} else {
		totalPages = alarmsCount / listCount
	}
	//获取当前页数
	curpage := this.Input().Get("CurPage")
	if curpage != "" {
		curPage, _ = strconv.ParseInt(curpage, 10, 0)
		beego.Debug("curPage=", curPage)
	}
	if curPage > totalPages {
		curPage = totalPages
	}
	//计算当前已显示条数
	curCount = listCount * (curPage - 1)

	//获取操作
	ope := this.Input().Get("op")
	switch ope {
	case "firstpage":
		curCount = 0
		curPage = 1
	case "prepage":
		if curPage > 1 {
			curCount -= listCount
			curPage -= 1
		}
	case "nextpage":
		if curPage < totalPages {
			curCount += listCount
			curPage += 1
		}
	case "lastpage":
		curCount = listCount * (totalPages - 1)
		curPage = totalPages
	case "delete":
		alarm := models.Alarm{
			Mac: mac,
		}
		models.DeleteAlarm(&alarm)

		path := "/alarm?CurUser=" + admin
		this.Redirect(path, 302)
		return
	}

	//从数据库中获取记录
	alarms, nums, ok := models.GetAlarms(listCount, curCount)
	if ok {
		this.Data["Alarms"] = alarms
		this.Data["AlarmsNum"] = alarmsCount
		this.Data["CurPage"] = curPage
		this.Data["ListCount"] = listCount
		this.Data["TotalPages"] = totalPages
	}
	if nums <= 0 {
		this.Data["NoInfo"] = "没有告警信息!"
	}
}

/*返回值
格式：{"code":0/-1/-2}
( 0) success
(-1) user/password error
(-2) other error
*/

func (this *AlarmController) Post() {
	ret := AlarmData{}
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
	alarm := models.Alarm{
		Time: time.Now(),
	}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &alarm)
	if err != nil {
		beego.Error(err)
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//将alarm信息写入数据库
	ok = models.AddAlarm(&alarm)
	if !ok {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}

	//通知客户端进行报警

	//返回成功
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	beego.Info(string(writeContent))
	return
}
