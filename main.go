package main

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"net"
	"tools-client/deploy"
	"tools-client/install"
)

func main() {
	if !install.IsInstall() {
		install.Install()
	}
	c := cron.New()
	c.AddFunc("*/3 * * * * ?", func() {
		go TaskLoop()
	}, "TaskLoop")
	c.Start()
	select {}
}

func TaskLoop() {
	conn, err := net.Dial("udp", "127.0.0.1:8002")
	defer conn.Close()
	if err != nil {
		panic("客户端启动失败-" + err.Error())
	}
	request := "001" + install.Cfg.String("zhimiao::DeployID")
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
