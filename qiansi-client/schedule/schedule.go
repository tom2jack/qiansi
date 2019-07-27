package schedule

import (
	"bytes"
	"encoding/gob"
	"github.com/jakecoffman/cron"
	"net"
	"qiansi/common/conf"
	"qiansi/common/models/udp_hook"
	"qiansi/qiansi-client/deploy"
)

var Cron *cron.Cron

func LoadSchedule() {
	Cron := cron.New()
	Cron.AddFunc("*/3 * * * * ?", TaskLoop, "TaskLoop")
	Cron.Start()
}

func TaskLoop() {
	conn, err := net.Dial("udp", "127.0.0.1:8001")
	defer conn.Close()
	if err != nil {
		return
	}
	request := "001" + conf.C.MustValue("zhimiao", "clientid")
	conn.Write([]byte(request))
	var r [255]byte
	n, _ := conn.Read(r[0:])
	resultBuf := bytes.NewBuffer(r[:n])
	var result = &udp_hook.Hook001DTO{}
	dec := gob.NewDecoder(resultBuf)
	dec.Decode(result)

	// 判断是否含有部署数据
	if result.Deploy != "" {
		deploy.Run(result.Deploy)
	}
}
