package main

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/conf"
	"gitee.com/zhimiao/qiansi/qiansi-client/install"
	"gitee.com/zhimiao/qiansi/qiansi-client/schedule"
	"os"
)

func main() {
	fmt.Print(`
          $$\                               $$\ 
          \__|                              \__|
 $$$$$$\  $$\  $$$$$$\  $$$$$$$\   $$$$$$$\ $$\ 
$$  __$$\ $$ | \____$$\ $$  __$$\ $$  _____|$$ |
$$ /  $$ |$$ | $$$$$$$ |$$ |  $$ |\$$$$$$\  $$ |
$$ |  $$ |$$ |$$  __$$ |$$ |  $$ | \____$$\ $$ |
\$$$$$$$ |$$ |\$$$$$$$ |$$ |  $$ |$$$$$$$  |$$ |
 \____$$ |\__| \_______|\__|  \__|\_______/ \__|
      $$ |                                      
      $$ |                                      
      \__|                                      
            
`)
	cfgFilePath := "config.ini"
	_, err := os.Stat(cfgFilePath)
	if err != nil {
		os.Create(cfgFilePath)
	}
	conf.C = conf.LoadConfig(cfgFilePath)
	// 判断安装
	if !install.IsInstall() {
		fmt.Println("配置文件读取失败，正在执行安装流程...")
		install.Install()
		os.Exit(0)
	}
	// 初始化计划任务
	schedule.LoadSchedule()
	fmt.Println(`千丝客户端启动完毕...`)
	// 阻塞
	select {}
}
