package install

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"qiansi/common/conf"
	"qiansi/qiansi-client/request"
	"strconv"
)

func IsInstall() bool {
	return conf.C.MustValue("zhimiao", "device") != "" && conf.C.MustValue("zhimiao", "clientid") != ""
}

func Install() bool {
	if conf.C.MustValue("zhimiao", "device") == "" {
		conf.C.SetValue("zhimiao", "device", uuid.NewV4().String())
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
		fmt.Print("请输入千丝平台注册ID:")
		num, err := fmt.Scanln(&UID)
		if err != nil || num == 0 || UID == "" {
			fmt.Println("数据读取失败，请重新输入！")
			continue
		} else {
			break
		}
	}
	fmt.Println("数据读取成功：" + UID)
	fmt.Println("当前机器唯一设备号：" + conf.C.MustValue("zhimiao", "device"))

	server, err := request.RegServer(UID, conf.C.MustValue("zhimiao", "device"))
	if err != nil {
		fmt.Println("注册失败了")
		return false
	}
	conf.C.SetValue("zhimiao", "uid", UID)
	conf.C.SetValue("zhimiao", "apisecret", server.ApiSecret)
	conf.C.SetValue("zhimiao", "clientid", strconv.Itoa(server.Id))
	err = conf.C.Save()
	if err != nil {
		fmt.Println("配置文件写入失败，请检查写入权限")
		return false
	}
	fmt.Println("Bind User Successful！")
	return true
}
