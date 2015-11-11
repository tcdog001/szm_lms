package controllers

import (
	"LMS/models"
	"github.com/astaxie/beego"
	"time"
)

/*
思路：用map来管理设备在线状态
（1）取命令、心跳请求，都会添加/维持当前在线状态
（2）启动goroutine，设定时间间隔，对map中的设备超时做检测，超时的设备修改状态为离线
*/

const (
	GC_INTERVAL      = 1  //(Minute)每隔一分钟清理一次超时设备
	TIMEOUT_INTERVAL = 10 //(Minute)超时时间为10分钟
)

type DeviceListener struct {
	State         int64
	LastAliveTime time.Time
}

var Listener map[string]DeviceListener

func init() {
	Listener = make(map[string]DeviceListener)

	go func() {
		for {
			for k, v := range Listener {
				beego.Debug("k=%s,v=%v", k, v)
				//判断是否超时
				if time.Now().Sub(v.LastAliveTime) >= time.Duration(TIMEOUT_INTERVAL)*time.Minute {
					//更新数据库中设备的状态为离线
					device := models.Deviceinfo{
						Mac:   k,
						State: 0,
					}
					models.UpdateDeviceStatus(&device)
					delete(Listener, k)
				}
			}
			beego.Debug("Listener GC runing...")
			time.Sleep(GC_INTERVAL * time.Minute)
		}
	}()
}
