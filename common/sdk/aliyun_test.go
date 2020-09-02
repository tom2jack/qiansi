package sdk

import (
	"github.com/zhi-miao/qiansi/common/config"
	"testing"
)

func TestAliyun_SendSmsVerify(t *testing.T) {
	config.LoadConfig("../../config.toml")
	NewAliyunSDK().SendSmsVerify("15061370322", "12")
}
