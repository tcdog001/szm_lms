package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//register all tables
	orm.RegisterModel(new(Userinfo), new(Admininfo), new(Deviceinfo), new(Historyinfo), new(Command), new(OperationRecord), new(Alarm), new(Script))
	//register mysql driver
	err := orm.RegisterDriver("mysql", orm.DR_MySQL)
	if err != nil {
		beego.Critical(err)
	}
	//register default database
	dbIp := beego.AppConfig.String("DbIp")
	dbPort := beego.AppConfig.String("DbPort")
	dbName := beego.AppConfig.String("DbName")
	dbUser := beego.AppConfig.String("DbUser")
	dbPassword := beego.AppConfig.String("DbPassword")

	dbUrl := dbUser + ":" + dbPassword + "@tcp(" + dbIp + ":" + dbPort + ")/" + dbName + "?charset=utf8&loc=Asia%2FShanghai"
	beego.Debug("dbUrl=", dbUrl)

	err = orm.RegisterDataBase("default", "mysql", dbUrl)
	if err != nil {
		beego.Critical(err)
	}
	//orm.RegisterDataBase("default", "mysql", "root:autelan@/lte_test?charset=utf8&&loc=Asia%2FShanghai")

	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)
}

func CheckDatabase() bool {
	o := orm.NewOrm()
	err := o.Using("default")
	if err != nil {
		return false
	} else {
		return true
	}
}
