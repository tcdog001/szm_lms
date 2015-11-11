package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Alarm struct {
	Id      int64     `orm:"auto";json:''-"`
	Mac     string    `json:"mac"`
	Content string    `json:"content"`
	Time    time.Time `json:''-"`
}

func (alarm *Alarm) TableName() string {
	return "alarm"
}

func AddAlarm(alarm *Alarm) bool {
	o := orm.NewOrm()
	_, err := o.Insert(alarm)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func DeleteAlarm(alarm *Alarm) bool {
	o := orm.NewOrm()
	var a Alarm
	err := o.QueryTable(alarm.TableName()).Filter("mac", alarm.Mac).One(&a)
	if err != nil {
		beego.Error(err)
		return false
	}
	_, err = o.Delete(&a)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func GetAlarms(start, offset int64) ([]*Alarm, int64, bool) {
	o := orm.NewOrm()
	var alarm Alarm
	//get all records
	alarms := make([]*Alarm, 0)
	num, err := o.QueryTable(alarm.TableName()).Limit(start, offset).All(&alarms)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return alarms, num, true
}

func GetAlarmCount() int64 {
	o := orm.NewOrm()
	var alarm Alarm
	cnt, _ := o.QueryTable(alarm.TableName()).Count()
	return cnt
}
