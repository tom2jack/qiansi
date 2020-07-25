package models

import (
	"fmt"
	"github.com/zhi-miao/qiansi/common"
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
