package main

import (
	"fmt"
	"os"
	"qiansi/common/conf"
	"qiansi/qiansi-client/install"
	"qiansi/qiansi-client/schedule"
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
	fmt.Print(`千丝客户端启动完毕...`)
	// 阻塞
	select {}
}
