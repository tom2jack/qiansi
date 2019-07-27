package schedule

import (
	"qiansi/common/conf"
	"qiansi/common/zmlog"
	"testing"
)

func TestTaskLoop(t *testing.T) {
	conf.C = conf.LoadConfig("./../../config.ini")
	zmlog.InitLog("./../../client.log")
	for i := 0; i < 70000; i++ {
		go TaskLoop()
		println(i)
	}
	select {}
}
