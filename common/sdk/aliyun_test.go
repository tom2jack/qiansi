package sdk

import (
	"github.com/zhi-miao/qiansi/common"
	"testing"
)

func TestAliyun_SendSmsVerify(t *testing.T) {
	common.Config.Init()
	NewAliyunSDK().SendSmsVerify("15061370322", "12")
}
