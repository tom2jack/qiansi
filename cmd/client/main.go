package main

import (
	"os"
	"qiansi/common/conf"
	"qiansi/common/zmlog"
	"qiansi/qiansi-client/install"
	"qiansi/qiansi-client/schedule"
)

func init() {
	conf.C = conf.LoadConfig("config.ini")
	zmlog.InitLog("client.log")
}

func main() {
	// 判断安装
	if !install.IsInstall() {
		install.Install()
		os.Exit(0)
	}
	// 初始化计划任务
	schedule.LoadSchedule()
	// 阻塞
	select {}
}
