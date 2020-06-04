package sdk

import (
	"gitee.com/zhimiao/qiansi/common"
	"testing"
)

func TestAliyun_SendSmsVerify(t *testing.T) {
	common.Config.Init()
	NewAliyunSDK().SendSmsVerify("15061370322", "12")
}
