package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Historyinfo struct {
	Id               int64     `orm:"auto";json:"-"`
	Mac              string    `json:"mac"`
	IpAddress        string    `json:"-"`
	IpLocation       string    `json:"-"`
	HostModel        string    `json:"hostModel"`
	HostSn           string    `json:"hostsn"`
	CpuModel         string    `json:"cpuModel"`
	CpuSn            string    `json:"cpuSN"`
	MemoryModel      string    `json:"memoryModel"`
	MemorySn         string    `json:"memorySN"`
	BoardSn          string    `json:"boardSN"`
	NetworkCardMac   string    `json:"networkCardMac"`
	DiskModel        string    `json:"diskModel"`
	DiskSn           string    `json:"diskSN"`
	LowfreModel      string    `json:"lowFreModel"`
	LowfreSn         string    `json:"lowFreSN"`
	HignfreModel     string    `json:"hignFreModel"`
	HignfreSn        string    `json:"hignFreSN"`
	GpsModel         string    `json:"gpsModel"`
	GpsSn            string    `json:"gpsSN"`
	Modelof3g        string    `json:"modelOf3g"`
	Snof3g           string    `json:"snOf3g"`
	Iccid            string    `json:"iccid"`
	HardVersion      string    `json:"hardVersion"`
	FirmwareVersion  string    `json:"firmwareVersion"`
	GatewayVersion   string    `json:"gateWayVersion"`
	ContentVersion   string    `json:"contentVersion"`
	RegistrationTime time.Time `orm:"auto_now_add;type(datetime)";json:"-"`
}

func (history *Historyinfo) TableName() string {
	return "historyinfo"
}

func RegisterHistory(historyinfo *Historyinfo) bool {
	o := orm.NewOrm()
	_, err := o.Insert(historyinfo)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}
