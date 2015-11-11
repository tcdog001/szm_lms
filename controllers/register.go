package controllers

import (
	"LMS/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/kayon/qqwry"
	//"strings"
	"time"
)

const (
	IPDAT      = "datas/qqwry.dat"
	ClientCert = "datas/root.crt"
)

type RegisterData struct {
	Code int64 `json:"code"`
}

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) Get() {
	this.TplNames = "home.html"
}

/*返回值
格式：{"code":0/-1/-2}
( 0) success
(-1) user/password error
(-2) other error
*/
func (this *RegisterController) Post() {
	//用户认证
	ret := RegisterData{}
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

	//从请求中获取客户端ip
	/*
		clientAddr := this.Ctx.Request.RemoteAddr
		index := strings.Index(clientAddr, ":")
		clientIp := clientAddr[:index]
		beego.Debug("client IP=", clientIp)
	*/
	clientIp := this.Ctx.Input.IP()
	//获取ip归属地
	iplocation := SelectIpLocation(clientIp)
	beego.Debug("client IpLocation=", iplocation)

	//获取请求的json数据
	beego.Info("request body=", string(this.Ctx.Input.RequestBody))
	deviceinfo := models.Deviceinfo{
		IpAddress:  clientIp,
		IpLocation: iplocation,
	}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &deviceinfo)
	if err != nil {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	} else {
		//注册设备到数据库
		if !models.RegisterDeivce(&deviceinfo) {
			ret.Code = -2
			writeContent, _ := json.Marshal(ret)
			this.Ctx.WriteString(string(writeContent))
			return
		}
	}
	beego.Info("insert table deviceinfo success!")

	//注册设备到历史记录表
	historyinfo := models.Historyinfo{
		IpAddress:  clientIp,
		IpLocation: iplocation,
	}
	json.Unmarshal(this.Ctx.Input.RequestBody, &historyinfo)
	if !models.RegisterHistory(&historyinfo) {
		ret.Code = -2
		writeContent, _ := json.Marshal(ret)
		this.Ctx.WriteString(string(writeContent))
		return
	}
	beego.Info("insert table historyinfo success!")

	//插入Listener中
	device := DeviceListener{
		State:         1,
		LastAliveTime: time.Now(),
	}
	Listener[deviceinfo.Mac] = device

	//返回注册成功
	ret.Code = 0
	writeContent, _ := json.Marshal(ret)
	this.Ctx.WriteString(string(writeContent))
	beego.Info(string(writeContent))
	return
}

func SelectIpLocation(ip string) string {
	qw := qqwry.New(IPDAT)
	var res qqwry.Result
	res = qw.Search(ip)
	beego.Debug("IP: %s\nBegin: %s\nEnd: %s\nCountry: %s\nArea: %s\n", res.IP, res.Begin, res.End, res.Country, res.Area)
	return res.Country + " " + res.Area
}
