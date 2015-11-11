package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type OperationRecord struct {
	Id       int64 `orm:"auto"`
	Operator string
	Mac      string
	Command  string `orm:"null"`
	Script   string `orm:"null"`
	Executed bool
	OpeTime  time.Time `orm:"auto_now_add;type(datetime)"`
	ExecTime time.Time `orm:"null"`
	Com      *Command  `orm:"reverse(one)"`
}

func (ope *OperationRecord) TableName() string {
	return "operation_record"
}

func AddOperationRecord(record *OperationRecord) bool {
	o := orm.NewOrm()
	_, err := o.Insert(record)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func UpdateOperationRecord(record *OperationRecord) bool {
	o := orm.NewOrm()
	/*
		var rec OperationRecord
		o.QueryTable(record.TableName()).Filter("mac", record.Mac).Filter("command", record.Command).Filter("script", record.Script).Filter("Executed", false).One(&rec)
		rec.Executed = record.Executed
		rec.ExecTime = record.ExecTime
		_, err := o.Update(&rec)
	*/
	var rec OperationRecord
	o.QueryTable(record.TableName()).Filter("id", record.Id).One(&rec)
	rec.Executed = record.Executed
	rec.ExecTime = record.ExecTime
	_, err := o.Update(&rec)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func GetRecords(start, offset int64) ([]*OperationRecord, int64, bool) {
	o := orm.NewOrm()
	var record OperationRecord
	//get all records
	records := make([]*OperationRecord, 0)
	num, err := o.QueryTable(record.TableName()).Limit(start, offset).All(&records)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return records, num, true
}

func GetRecordCount() int64 {
	o := orm.NewOrm()
	var record OperationRecord
	cnt, _ := o.QueryTable(record.TableName()).Count()
	return cnt
}
