package controllers

/*
import (
	"LMS/models"
	"os"
	"time"
)

//定时(每天凌晨三点钟)扫描script表中，发现状态为downloaded的script，则删除

var firsttime = true


func init() {
	go func() {
		for {
			if firsttime {
				//计算此刻到凌晨3点所需时间
				d := 24 - time.Now().Hour() + 3
				time.Sleep(time.Duration(d) * time.Hour)
				firsttime = false
			} else {
				//一天之后
				d, _ := time.ParseDuration("24h")
				time.Sleep(d)
			}

			//time.Sleep(10 * time.Second)
			scripts, ok := models.GetAllScripts()
			if ok {
				for _, script := range scripts {
					if script.Downloaded {
						//删除硬盘上对应的脚本
						os.Remove(ScriptFolder + script.FilePath)
						//删除script表对应项
						models.DeleteScript(script)
					}
				}
			}
		}
	}()
}
*/
