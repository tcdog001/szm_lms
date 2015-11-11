package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"
const TimeConvert = "0001-01-01 00:00:00"

type Deviceinfo struct {
	Mac                   string    `orm:"pk";json:"mac"`
	IpAddress             string    `json:"-"`
	IpLocation            string    `json:"-"`
	HostModel             string    `json:"hostModel"`
	HostSn                string    `json:"hostsn"`
	CpuModel              string    `json:"cpuModel"`
	CpuSn                 string    `json:"cpuSN"`
	MemoryModel           string    `json:"memoryModel"`
	MemorySn              string    `json:"memorySN"`
	BoardSn               string    `json:"boardSN"`
	NetworkCardMac        string    `json:"networkCardMac"`
	DiskModel             string    `json:"diskModel"`
	DiskSn                string    `json:"diskSN"`
	LowfreModel           string    `json:"lowFreModel"`
	LowfreSn              string    `json:"lowFreSN"`
	HignfreModel          string    `json:"hignFreModel"`
	HignfreSn             string    `json:"hignFreSN"`
	GpsModel              string    `json:"gpsModel"`
	GpsSn                 string    `json:"gpsSN"`
	Modelof3g             string    `json:"modelOf3g"`
	Snof3g                string    `json:"snOf3g"`
	Iccid                 string    `json:"iccid"`
	HardVersion           string    `json:"hardVersion"`
	FirmwareVersion       string    `json:"firmwareVersion"`
	GatewayVersion        string    `json:"gateWayVersion"`
	ContentVersion        string    `json:"contentVersion"`
	FirstRegistrationTime time.Time `orm:"auto_now_add;type(datetime)";json:"-"`
	LastRegistrationTime  time.Time `json:"-"`
	State                 int       `orm:"default(0)";json:"-"`
	LastKeepaliveTime     time.Time `orm:"auto_now;type(datetime)";json:"-"`
}

func (device *Deviceinfo) TableName() string {
	return "deviceinfo"
}

func RegisterDeivce(device *Deviceinfo) bool {
	o := orm.NewOrm()
	//device.LastKeepaliveTime = time.Now()
	device.LastRegistrationTime = time.Now()
	device.State = 1
	beego.Debug("RegisterDeivce table=", device.TableName())
	//查找对应的mac地址是否存在，存在的话要求状态为离线
	exist := o.QueryTable(device.TableName()).Filter("mac", device.Mac).Filter("state", 0).Exist()
	if exist {
		//设备存在，则更新设备信息
		return UpdateDevice(device)
	} else {
		//设备不存在，则插入设备信息
		_, err := o.Insert(device)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func UpdateDevice(device *Deviceinfo) bool {
	beego.Debug("UpdateDevice table=", device.TableName())
	o := orm.NewOrm()

	var dev Deviceinfo
	err := o.QueryTable(device.TableName()).Filter("mac", device.Mac).One(&dev)
	if err != nil {
		return false
	} else {
		dev.State = device.State
		dev.LastKeepaliveTime = device.LastKeepaliveTime
		dev.LastRegistrationTime = device.LastRegistrationTime
		dev.IpAddress = device.IpAddress
		dev.IpLocation = device.IpLocation

		beego.Debug("UpdateDevice clientIp =", dev.IpAddress)

		_, err := o.Update(&dev)
		if err != nil {
			beego.Error(err)
			return false
		}
		return true
	}
}

func UpdateDeviceStatus(device *Deviceinfo) bool {
	o := orm.NewOrm()

	var dev Deviceinfo
	err := o.QueryTable(device.TableName()).Filter("mac", device.Mac).One(&dev)
	if err != nil {
		beego.Error(err)
		return false
	}
	dev.State = device.State
	dev.LastKeepaliveTime = device.LastKeepaliveTime
	_, err = o.Update(&dev)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func GetDevices(start, offset int64) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	var device Deviceinfo
	//get all devices
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable(device.TableName()).Limit(start, offset).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetDevicesCount() int64 {
	o := orm.NewOrm()
	var device Deviceinfo
	cnt, _ := o.QueryTable(device.TableName()).Count()
	return cnt
}

func GetSearchDevices(start, offset int64, device *Deviceinfo) ([]*Deviceinfo, int64, bool) {
	o := orm.NewOrm()
	cond := orm.NewCondition()

	if device.Mac != "" {
		cond = cond.And("mac__icontains", device.Mac)
	}
	if device.IpAddress != "" {
		cond = cond.And("ip_address__icontains", device.IpAddress)
	}
	if device.IpLocation != "" {
		cond = cond.And("ip_location__icontains", device.IpLocation)
	}
	if device.DiskSn != "" {
		cond = cond.And("disk_sn__icontains", device.HostSn)
	}
	if device.FirmwareVersion != "" {
		cond = cond.And("firmware_version__icontains", device.FirmwareVersion)
	}
	if device.State != -1 {
		cond = cond.And("state__icontains", device.State)
	}

	/*
		t, _ := time.Parse(TimeFormat, TimeConvert)
		beego.Debug("t=", t)
		beego.Debug("FirstRegistrationTime=", device.FirstRegistrationTime)
		beego.Debug(t == device.FirstRegistrationTime)
		if device.FirstRegistrationTime != t {
			cond = cond.And("first_registration_time__lte", device.FirstRegistrationTime)
		}
		if device.LastRegistrationTime != t {
			cond = cond.And("last_registration_time__gte", device.LastRegistrationTime)
		}
	*/
	//get all devices
	devices := make([]*Deviceinfo, 0)
	num, err := o.QueryTable(device.TableName()).Limit(start, offset).SetCond(cond).All(&devices)
	if err != nil {
		beego.Error(err)
		return nil, 0, false
	}
	return devices, num, true
}

func GetSearchDevicesCount(device *Deviceinfo) int64 {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	if device.Mac != "" {
		cond = cond.And("mac__icontains", device.Mac)
	}
	if device.IpAddress != "" {
		cond = cond.And("ip_address__icontains", device.IpAddress)
	}
	if device.IpLocation != "" {
		cond = cond.And("ip_location__icontains", device.IpLocation)
	}
	if device.DiskSn != "" {
		cond = cond.And("disk_sn__icontains", device.HostSn)
	}
	if device.FirmwareVersion != "" {
		cond = cond.And("firmware_version__icontains", device.FirmwareVersion)
	}
	if device.State != -1 {
		cond = cond.And("state__icontains", device.State)
	}
	/*
		t, _ := time.Parse(TimeFormat, TimeConvert)
		beego.Debug("t=", t)
		if device.FirstRegistrationTime != t {
			cond = cond.And("first_registration_time__lte", device.FirstRegistrationTime)
		}
		if device.LastRegistrationTime != t {
			cond = cond.And("last_registration_time__gte", device.LastRegistrationTime)
		}
	*/
	cnt, _ := o.QueryTable(device.TableName()).SetCond(cond).Count()

	beego.Debug("search devices count=", cnt)
	return cnt
}
