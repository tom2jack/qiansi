package install

import (
	"fmt"
)

func Install() bool {
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
	return true
}
