package conf

import (
	"github.com/Unknwon/goconfig"
	"os"
)

var (
	S   *ZMCfg
	C   *ZMCfg
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
		file, err := os.Create(path)
		if err != nil {
			panic("文件无法写入")
		}
		file.Close()
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
