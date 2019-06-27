package common

import (
	"github.com/astaxie/beego/config"
	"os"
)

var (
	Cfg config.Configer = InitConfig()
)

func GetConfigPath() string {
	return "config.ini"
}

//InitConfig 初始化配置
func InitConfig() config.Configer {
	_, err := os.Stat(GetConfigPath())
	if err != nil {
		file, err := os.Create(GetConfigPath())
		if err != nil {
			panic("文件无法写入")
		}
		file.Close()
	}
	cfg, err := config.NewConfig("ini", GetConfigPath())
	if err != nil {
		panic("配置文件读取错误")
	}
	return cfg
}
