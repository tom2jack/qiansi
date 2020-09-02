package config

import (
	"github.com/jinzhu/configor"
)

type configStruct struct {
	App struct {
		PageSize  int
		JwtSecret string
	}
	Server struct {
		ApiHost      string
		APIListen    string
		UDPListen    string
		ReadTimeOut  int
		WriteTimeOut int
	}
	Mysql struct {
		Host        string
		User        string
		Password    string
		Database    string
		TablePrefix string
	}
	InfluxDB struct {
		Host  string
		Token string
		Org   string
	}
	Redis struct {
		Host string
		Auth string
		DB   int
	}
	Aliyun struct {
		AccessKey    string
		AccessSecret string
		RegionId     string
		SmsConfig    struct {
			SignName     string
			TemplateCode string
			RegionId     string
		}
	}
	MQTT struct {
		Broker   string
		Username string
		Password string
		ClientID string
	}
}

var cfg = &configStruct{}

// Init 初始化配置
func LoadConfig(filePath string) error {
	return configor.Load(cfg, filePath)
}

// GetConfig 获取配置
func GetConfig() *configStruct {
	return cfg
}

// ENV 获取当前配置场景
func ENV() string {
	return configor.ENV()
}
