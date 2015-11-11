package controllers

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
)

type ShowScriptController struct {
	beego.Controller
}

func (this *ShowScriptController) Get() {
	//用户身份认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}
	//获取脚本名称
	scriptname := this.Input().Get("script")
	beego.Debug("get script name=", scriptname)
	//读取脚本中的内容并返回客户端
	scriptfile := ScriptFolder + scriptname
	this.Ctx.WriteString(readfile(scriptfile))
}

func readfile(path string) string {
	fd, err := os.Open(path)
	if err != nil {
		beego.Error(err)
	}
	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		beego.Error(err)
	}
	return string(content)
}

func (this *ShowScriptController) Post() {
	//用户身份认证
	session := this.GetSession("Admin")
	if session == nil {
		beego.Trace("session verify failed!")
		this.Redirect("/", 302)
		return
	}

}
