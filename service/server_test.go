package service

import (
	"fmt"
	"github.com/zhi-miao/qiansi/common/config"
	"testing"
)

func TestGetServerService(t *testing.T) {
	config.LoadConfig("config.dev.toml")
	source, err := GetServerService().GetClientSource("linux", "amd64")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(*source)
}
