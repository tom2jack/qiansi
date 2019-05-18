package install

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Install() bool {
	//绑定用户
	if !binUser() {
		return false
	}

	return true
}

func binUser() bool {
	fmt.Print("请输入纸喵运维平台注册ID:")
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("\r\n" + err.Error())
		return false
	}
	str = strings.TrimSpace(str)
	if str == "" {
		fmt.Println("\r\n非法输入，已退出")
	}
	fmt.Println(str)
	return true
}
