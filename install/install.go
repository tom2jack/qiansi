package install

import (
	"fmt"
	"github.com/astaxie/beego/config"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	CfgFilePath string          = "config.ini"
	Cfg         config.Configer = InitConfig()
)

func Install() bool {
	//绑定用户
	if !binUser() {
		return false
	}
	return true
}

//InitConfig 初始化配置
func InitConfig() config.Configer {
	_, err := os.Stat(CfgFilePath)
	if err != nil {
		file, err := os.Create(CfgFilePath)
		if err != nil {
			panic("文件无法写入")
		}
		file.Close()
	}
	cfg, err := config.NewConfig("ini", CfgFilePath)
	if err != nil {
		panic("配置文件读取错误")
	}
	if cfg.String("zhimiao::device") == "" {
		cfg.Set("zhimiao::device", uuid.NewV4().String())
		cfg.SaveConfigFile(CfgFilePath)
	}
	return cfg
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
	fmt.Println("当前机器唯一设备号：" + Cfg.String("zhimiao::device"))
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
