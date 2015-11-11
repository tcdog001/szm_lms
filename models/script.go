package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Script struct {
	Id           int64  `orm:"auto"`
	Mac          string `json:"mac"`
	FilePath     string
	Downloaded   bool
	ScriptTime   time.Time `orm:"auto_now;type(datetime)"`
	DowonlodTime time.Time `orm:"null"`
}

func (script *Script) TableName() string {
	return "script"
}

func AddScript(script *Script) bool {
	o := orm.NewOrm()
	_, err := o.Insert(script)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func GetScript(script *Script) (*Script, bool) {
	o := orm.NewOrm()
	var sc Script
	//获取对应mac下的脚本，并且取时间为最早时间
	err := o.QueryTable(script.TableName()).Filter("mac", script.Mac).Filter("downloaded", script.Downloaded).OrderBy("script_time").One(&sc)
	if err != nil {
		beego.Error(err)
		return nil, false
	}
	return &sc, true
}

func GetAllScripts() ([]*Script, bool) {
	o := orm.NewOrm()
	var script Script

	scripts := make([]*Script, 0)
	_, err := o.QueryTable(script.TableName()).Filter("downloaded", true).All(&scripts)
	if err != nil {
		beego.Error(err)
		return nil, false
	}
	return scripts, true
}

func UpdateScript(script *Script) bool {
	o := orm.NewOrm()
	_, err := o.Update(script)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func DeleteScript(script *Script) bool {
	o := orm.NewOrm()
	_, err := o.Delete(script)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
