package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Userinfo struct {
	Uid      int64 `orm:"auto"`
	Username string
	Password string
	Created  time.Time
}

func (user *Userinfo) TableName() string {
	return "userinfo"
}

func CheckAccount(user *Userinfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable(user.TableName()).Filter("username__exact", user.Username).Filter("password__exact", user.Password).Exist()
	return exist
}
