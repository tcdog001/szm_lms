package routers

import (
	"LMS/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
)

func init() {
	toolbox.AddHealthCheck("database", &controllers.DatabaseCheck{})
	//网管后台路由
	beego.Router("/", &controllers.LoginController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/home", &controllers.HomeController{})
	beego.Router("/deviceinfo", &controllers.DeviceinfoController{})
	beego.Router("/search", &controllers.SearchController{})
	beego.Router("/command", &controllers.CommandController{})
	beego.Router("/operecord", &controllers.OperecordController{})
	beego.Router("/script", &controllers.ScriptController{})
	beego.Router("/showscript", &controllers.ShowScriptController{})
	beego.Router("/alarm", &controllers.AlarmController{})

	//与设备通信路由
	beego.Router("/LMS/lte/lteRegister.do", &controllers.RegisterController{})
	beego.Router("/LMS/lte/lteUnRegister.do", &controllers.UnRegisterController{})
	beego.Router("/LMS/lte/lteCommand.do", &controllers.GetCommandController{})
	beego.Router("/LMS/lte/lteHeartBeat.do", &controllers.HeartBeatController{})
	beego.Router("/LMS/lte/lteScript.do", &controllers.GetScriptController{})
	beego.Router("/LMS/lte/lteAlarm.do", &controllers.AlarmController{})
}
