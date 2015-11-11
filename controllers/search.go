package controllers

import (
	"LMS/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

type SearchController struct {
	beego.Controller
}

const TimeFormat = "2006-01-02 15:04:05"

func (this *SearchController) Get() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	this.TplNames = "search.html"

	CurUser := this.Input().Get("CurUser")
	this.Data["CurUser"] = CurUser

	//获取搜索条件内容
	mac := this.Input().Get("Mac")
	ip := this.Input().Get("Ip")
	dsn := this.Input().Get("Dsn")
	lc := this.Input().Get("Lc")
	fr := this.Input().Get("Fr")
	lr := this.Input().Get("Lr")
	fw := this.Input().Get("Fw")
	ds := this.Input().Get("Ds")
	beego.Debug("device=", mac, ip, lc, dsn, fr, lr, fw, ds)

	this.Data["Mac"] = mac
	this.Data["Ip"] = ip
	this.Data["Dsn"] = dsn
	this.Data["Lc"] = lc
	this.Data["Fr"] = fr
	this.Data["Lr"] = lr
	this.Data["Fw"] = fw
	this.Data["Ds"] = ds

	if mac == "" && ip == "" && dsn == "" && lc == "" && fr == "" && lr == "" && fw == "" && ds == "" {
		fmt.Println("searchCondition is null")
		this.Data["Show"] = false
	} else {
		beego.Debug("searchCondition is not null")
		this.Data["Show"] = true

		//重置参数
		devicesCount = 0
		totalPages = 0
		curPage = 1
		listCount = 10
		curCount = 0

		//fr to time.Time
		ft, err := time.Parse(TimeFormat, fr)
		if err != nil {
			beego.Error(err)
		}
		//lr to time.Time
		lt, err := time.Parse(TimeFormat, lr)
		if err != nil {
			beego.Error(err)
		}
		//ds to int
		state, _ := strconv.Atoi(ds)

		device := models.Deviceinfo{
			Mac:                   mac,
			IpAddress:             ip,
			IpLocation:            lc,
			HostSn:                dsn,
			FirstRegistrationTime: ft,
			LastRegistrationTime:  lt,
			FirmwareVersion:       fw,
			State:                 state,
		}

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
			devicesCount = models.GetSearchDevicesCount(&device)
		}
		//计算总页数
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
		}
		if curPage > totalPages {
			curPage = totalPages
		}
		//获取当前已显示数
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

		//根据条件，从数据库获取设备
		devices, nums, ok := models.GetSearchDevices(listCount, curCount, &device)
		if ok {
			this.Data["Devices"] = devices
			this.Data["DevicesNum"] = devicesCount
			this.Data["CurPage"] = curPage
			this.Data["ListCount"] = listCount
			this.Data["TotalPages"] = totalPages
		}
		if nums <= 0 {
			this.Data["NoInfo"] = "符合条件设备不存在!"
		}
	}
}

func (this *SearchController) Post() {
	//session认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	CurUser := this.Input().Get("CurUser")

	searchCondition := this.Input().Get("searchCondition")
	beego.Debug("searchCondition input=", searchCondition)
	device := ParseSearchContent(searchCondition)
	beego.Debug("device=", *device)

	FrTime := Substr(device.FirstRegistrationTime.String(), 0, 19)
	LrTime := Substr(device.LastRegistrationTime.String(), 0, 19)
	beego.Debug("FrTime=", FrTime)
	beego.Debug("LrTime=", LrTime)

	deviceInfo := "Mac=" + device.Mac +
		"&Ip=" + device.IpAddress +
		"&Dsn=" + device.HostSn +
		"&Lc=" + device.IpLocation +
		"&Fr=" + FrTime +
		"&Lr=" + LrTime +
		"&Fw=" + device.FirmwareVersion +
		"&Ds=" + strconv.Itoa(device.State) +
		"&CurUser=" + CurUser
	//this.Data["Deviceinfo"] = deviceInfo

	search := "search?" + deviceInfo
	beego.Debug("search url=", search)
	this.Redirect(search, 302)
	return
}

//解析搜索条件内容
func ParseSearchContent(search string) *models.Deviceinfo {
	beego.Debug("searchCondition=", search)

	searchInfo := make(map[string]string)
	searchs := strings.Split(search, "&&")
	for _, item := range searchs {
		if item == "" {
			continue
		} else {
			info := strings.Split(item, "=")
			searchInfo[info[0]] = info[1]
		}
	}
	beego.Debug("searchInfo=", searchInfo)

	ds_flag := false

	var device models.Deviceinfo
	for key, value := range searchInfo {
		//去掉所有空格
		key = strings.Replace(key, " ", "", -1)
		value = strings.Replace(value, " ", "", -1)
		beego.Debug("key=", key, "value=", value)

		switch key {
		case "mac":
			device.Mac = value
		case "ip":
			device.IpAddress = value
		case "dsn":
			device.HostSn = value
		case "lc":
			device.IpLocation = value
		case "fr":
			if value != "" {
				t, err := time.Parse(TimeFormat, value)
				if err != nil {
					beego.Error(err)
				}
				device.FirstRegistrationTime = t
			}
		case "lr":
			if value != "" {
				t, err := time.Parse(TimeFormat, value)
				if err != nil {
					beego.Error(err)
				}
				device.LastRegistrationTime = t
			}
		case "fw":
			device.FirmwareVersion = value
		case "ds":
			ds_flag = true
			device.State, _ = strconv.Atoi(value)
		default:
		}
	}
	if !ds_flag {
		device.State = -1
	}
	return &device
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
