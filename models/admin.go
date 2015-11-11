package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Admininfo struct {
	Uid           int64     `orm:"auto"`
	Username      string    `json:"uname"`
	Password      string    `json:"pwd"`
	CreatedTime   time.Time `json:"-"`
	LastLoginTime time.Time `orm:"null;auto_now;type(datetime)";json:"lastlogin"`
}

func (admin *Admininfo) TableName() string {
	return "admininfo"
}

func CheckAdmin(admin *Admininfo) bool {
	o := orm.NewOrm()
	exist := o.QueryTable(admin.TableName()).Filter("username__exact", admin.Username).Filter("password__exact", admin.Password).Exist()
	return exist
}

func UpdateAdminStatus(admin *Admininfo) bool {
	o := orm.NewOrm()
	var u Admininfo
	err := o.QueryTable(admin.TableName()).Filter("username", admin.Username).One(&u)
	if err != nil {
		beego.Error(err)
		return false
	}
	u.LastLoginTime = time.Now()
	_, err = o.Update(&u)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
