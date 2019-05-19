package install

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Install() bool {
	cfg()
	//fmt.Println(os.Hostname())
	//fmt.Println(common.GetLocalMac())
	//绑定用户
	if !binUser() {
		return false
	}
	return true
}

func cfg() {
	_, err := os.Stat("config.ini")
	if err != nil {
		file, err := os.Create("config.ini")
		if err != nil {
			panic("文件无法写入")
		}
		file.Close()
	}
	cfg, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		panic("配置文件读取错误")
	}
	cfg.Set("zhimiao::device", uuid.NewV4().String())
	cfg.SaveConfigFile("config.ini")

}

func binUser() bool {
	var UID string
	for {
		fmt.Print("请输入纸喵运维平台注册ID:")
		num, err := fmt.Scanln(&UID)
		if err != nil || num == 0 || UID == "" {
			fmt.Println("数据读取失败，请重新输入！")
			continue
		} else {
			break
		}
	}
	fmt.Println("数据读取成功：" + UID)
	return true
}

func RequestJob(UID string, HostName string) {
	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	url := "http://localhost:1305/index.php?c=api&a=regServer&uid=" + UID + "&hostname=" + HostName
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != "2" {
		return
	}

}
