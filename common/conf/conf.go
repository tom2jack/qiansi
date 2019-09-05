package conf

import (
	"os"
	"github.com/Unknwon/goconfig"
)

var (
	S   *ZMCfg = LoadConfig("assets/config/server.ini")
	C   *ZMCfg = LoadConfig("config.ini")
	err error
)

type ZMCfg struct {
	goconfig.ConfigFile
	FilePath string
}

func (cfg *ZMCfg) Save() error {
	return goconfig.SaveConfigFile(&cfg.ConfigFile, cfg.FilePath)
}

func LoadConfig(path string) *ZMCfg {
	_, err = os.Stat(path)
	if err != nil {
		return &ZMCfg{
			ConfigFile: goconfig.ConfigFile{},
			FilePath:   "",
		}
	}
	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		panic("app配置文件读取异常")
	}
	return &ZMCfg{
		FilePath:   path,
		ConfigFile: *cfg,
	}
}
