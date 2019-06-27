package main

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"net"
	"os"
	"tools-client/common"
	"tools-client/deploy"
	"tools-client/install"
)

func main() {
	// 判断安装
	if !install.IsInstall() {
		install.Install()
		os.Exit(0)
	}
	// 主程启动，定时扫描任务
	c := cron.New()
	c.AddFunc("*/3 * * * * ?", TaskLoop, "TaskLoop")
	c.Start()
	select {}
}

func TaskLoop() {
	conn, err := net.Dial("udp", "127.0.0.1:8002")
	defer conn.Close()
	if err != nil {
		panic("客户端启动失败-" + err.Error())
	}
	request := "001" + common.Cfg.String("zhimiao::clientid")
	conn.Write([]byte(request))
	var result [1]byte
	conn.Read(result[0:])
	data := string(result[0:1])
	fmt.Println("request:", request)
	fmt.Println("msg is", data)

	switch data {
	case "1":
		go deploy.Run()
	default:
		fmt.Println("Task loop miss:", data)
	}
}
