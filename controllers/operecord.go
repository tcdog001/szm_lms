package controllers

import (
	"LMS/models"
	"github.com/astaxie/beego"
	"strconv"
)

type OperecordController struct {
	beego.Controller
}

var recordsCount int64

func (this *OperecordController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "opeaterecord.html"
	CurUser := this.Input().Get("CurUser")
	beego.Debug("admin=", CurUser)
	this.Data["CurUser"] = CurUser

	//重置参数
	recordsCount = 0
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
	recordscount := this.Input().Get("RecordsCount")
	if recordscount != "" {
		recordsCount, _ = strconv.ParseInt(recordscount, 10, 0)
		beego.Debug("recordsCount=", recordsCount)
	} else {
		recordsCount = models.GetRecordCount()
	}
	//计算总页数
	if recordsCount%listCount > 0 {
		totalPages = recordsCount/listCount + 1
	} else {
		totalPages = recordsCount / listCount
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
	}

	//从数据库中获取记录
	records, nums, ok := models.GetRecords(listCount, curCount)
	if ok {
		this.Data["Records"] = records
		this.Data["RecordsNum"] = recordsCount
		this.Data["CurPage"] = curPage
		this.Data["ListCount"] = listCount
		this.Data["TotalPages"] = totalPages
	}
	if nums <= 0 {
		this.Data["NoInfo"] = "没有操作记录!"
	}
}

func (this *OperecordController) Post() {
	//sesson认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.Redirect("/", 302)
	return
}
