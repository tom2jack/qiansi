package conf

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"path/filepath"
)

var (
	S    *ZMCfg
	C    *ZMCfg
	ROOT string
	err  error
)

const (
	API_URL  = "http://qiansi.zhimiao.org"
	UDP_HOST = "qiansi.zhimiao.org:1315"
)

type ZMCfg struct {
	goconfig.ConfigFile
	FilePath string
}

func init() {
	ROOT, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	S = LoadConfig("assets/config/server.ini")
}

func (cfg *ZMCfg) Save() error {
	return goconfig.SaveConfigFile(&cfg.ConfigFile, cfg.FilePath)
}

func LoadConfig(path string) *ZMCfg {
	_, err = os.Stat(path)
	if err != nil {
		return nil
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
