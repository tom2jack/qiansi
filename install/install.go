package install

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"tools-client/common"
	"tools-client/request"
)

var (
	cfg = common.Cfg
)

func IsInstall() bool {
	return cfg.String("zhimiao::device") != "" && cfg.String("zhimiao::clientid") != ""
}

func Install() bool {
	if cfg.String("zhimiao::device") == "" {
		cfg.Set("zhimiao::device", uuid.NewV4().String())
	}
	//绑定用户
	if !binUser() {
		return false
	}
	return true
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
	fmt.Println("当前机器唯一设备号：" + cfg.String("zhimiao::device"))

	server, err := request.RegServer(UID, cfg.String("zhimiao::device"))
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg.Set("zhimiao::uid", UID)
	cfg.Set("zhimiao::apisecret", server.ApiSecret)
	cfg.Set("zhimiao::clientid", strconv.Itoa(server.Id))
	cfg.SaveConfigFile(common.GetConfigPath())
	fmt.Println("Bind User Successful！")
	return true
}
