package controllers

import (
	"LMS/models"
	"github.com/astaxie/beego"
	"os"
	"strconv"
	"time"
)

type ScriptController struct {
	beego.Controller
}

var ScriptFolder = "script_download/"

func init() {
	if !Exist(ScriptFolder) {
		os.Mkdir(ScriptFolder, 0)
	}
}

func (this *ScriptController) Get() {
	//用户身份认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	this.TplNames = "script.html"
	mac := this.Input().Get("mac")
	admin := this.Input().Get("CurUser")
	this.Data["Mac"] = mac
	this.Data["CurUser"] = admin
}

/*返回值
格式：{"code":0/-1/-2}
( 0) success
(-1) user/password error
(-2) other error
*/

func (this *ScriptController) Post() {
	//用户身份认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

	//获取post信息
	mac := this.Input().Get("Mac")
	admin := this.Input().Get("CurUser")
	beego.Debug(mac)
	beego.Debug(admin)
	//获取上传脚本
	filename := strconv.FormatInt(int64(time.Now().UnixNano()), 10) + ".sh"
	filepath := ScriptFolder + filename
	beego.Debug("saveFileName=", filepath)

	this.SaveToFile("file", filepath)

	//将script信息写入数据库
	script := models.Script{
		Mac:        mac,
		FilePath:   filename,
		Downloaded: false,
	}
	ok := models.AddScript(&script)
	if !ok {
		url := this.Ctx.Request.URL
		beego.Info(url)
		path := "/deviceinfo?mac=" + mac + "&CurUser=" + admin
		this.Redirect(path, 302)
		return
	}

	//数据库中添加管理员操作记录
	record := models.OperationRecord{
		Operator: admin,
		Mac:      mac,
		Script:   filename,
		Executed: false,
	}
	models.AddOperationRecord(&record)

	//返回
	path := "/deviceinfo?mac=" + mac + "&CurUser=" + admin
	this.Redirect(path, 302)
	return
}

// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
