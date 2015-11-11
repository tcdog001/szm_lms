package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Command struct {
	Id        int64 `orm:"auto"`
	Mac       string
	Command   string
	Executed  bool
	Created   time.Time        `orm:"auto_now_add;type(datetime)"`
	Operecord *OperationRecord `orm:"rel(one)"`
}

func (com *Command) TableName() string {
	return "command"
}

func AddDeviceCommand(command *Command) bool {
	o := orm.NewOrm()
	_, err := o.Insert(command)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func UpdateDeviceCommand(command *Command) bool {
	o := orm.NewOrm()

	var com Command
	o.QueryTable(command.TableName()).Filter("id", command.Id).One(&com)
	beego.Debug(com)
	if !command.Executed {
		com.Executed = true
	} else {
		com.Executed = command.Executed
	}
	_, err := o.Update(&com)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func CheckCommandExist(command *Command) bool {
	o := orm.NewOrm()
	exist := o.QueryTable(command.TableName()).Filter("mac", command.Mac).Filter("executed", command.Executed).Exist()
	return exist
}

func GetCommand(command *Command) (bool, *Command) {
	o := orm.NewOrm()
	var com Command
	var coms []Command
	coms = make([]Command, 0)
	num, err := o.QueryTable(command.TableName()).Filter("mac", command.Mac).Filter("executed", false).OrderBy("created").All(&coms)
	beego.Debug(num, coms, err)
	if num <= 0 {
		return false, &com
	}
	for k, v := range coms {
		//只取第一个元素
		if k == 0 {
			com = v
			break
		}
	}
	return true, &com
}
