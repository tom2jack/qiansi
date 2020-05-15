package models

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common"
	"os"
	"testing"
)

func TestDeployInfo(t *testing.T) {
	os.Setenv("CONFIGOR_ENV", "dev")
	common.Config.Init()
	Start()
	info := DeployInfo(34)
	fmt.Println(info)
}
