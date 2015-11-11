package controllers

import (
	"LMS/models"
	"github.com/astaxie/beego"
	"strconv"
)

var (
	devicesCount int64 = 0  //默认设备数
	totalPages   int64 = 0  //默认总页数
	curPage      int64 = 1  //默认当前页数
	listCount    int64 = 10 //默认每页显示条数
	curCount     int64 = 0  //默认已经显示的设备数
)

type DeviceinfoController struct {
	beego.Controller
}

func (this *DeviceinfoController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "deviceinfo.html"
	CurUser := this.Input().Get("CurUser")
	beego.Debug("admin=", CurUser)
	this.Data["CurUser"] = CurUser

	//重置参数
	devicesCount = 0
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

	//获取设备总数
	devicescount := this.Input().Get("DevicesCount")
	if devicescount != "" {
		devicesCount, _ = strconv.ParseInt(devicescount, 10, 0)
		beego.Debug("devicesCount=", devicesCount)
	} else {
		devicesCount = models.GetDevicesCount()
	}

	//计算出总页数
	if devicesCount%listCount > 0 {
		totalPages = devicesCount/listCount + 1
	} else {
		totalPages = devicesCount / listCount
	}

	//获取当前页数
	curpage := this.Input().Get("CurPage")
	if curpage != "" {
		curPage, _ = strconv.ParseInt(curpage, 10, 0)
		beego.Debug("curPage=", curPage)
	} else {
		curPage = 1
	}
	if curPage > totalPages {
		curPage = totalPages
	}

	//计算出当前总条数
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
	}

	devices, nums, ok := models.GetDevices(listCount, curCount)
	if ok {
		this.Data["Devices"] = devices
		this.Data["DevicesNum"] = devicesCount
		this.Data["CurPage"] = curPage
		this.Data["ListCount"] = listCount
		this.Data["TotalPages"] = totalPages
	}
	if nums <= 0 {
		this.Data["CurPage"] = 0
		this.Data["NoInfo"] = "没有注册设备!"
	}
}

func (this *DeviceinfoController) Post() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.Redirect("/", 302)
	return
}
