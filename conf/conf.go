package conf

import (
	"github.com/Unknwon/goconfig"
)

var (
	App *goconfig.ConfigFile
)

func LoadConfig() {
	App = AppConfig()
}

func AppConfig() *goconfig.ConfigFile {
	cfg, err := goconfig.LoadConfigFile("conf/app.ini")
	if err != nil {
		panic("app配置文件读取异常")
	}
	return cfg
}
