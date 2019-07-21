package request

import (
	"qiansi/common/conf"
	"testing"
)

func TestLogPush(t *testing.T) {
	conf.C = conf.LoadConfig("../../config.ini")
	LogPush("ballll")
}
