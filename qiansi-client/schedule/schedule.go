package schedule

import (
	"bytes"
	"encoding/gob"
	"gitee.com/zhimiao/qiansi/common/conf"
	"gitee.com/zhimiao/qiansi/common/dto"
	"gitee.com/zhimiao/qiansi/qiansi-client/deploy"
	"github.com/jakecoffman/cron"
	"net"
	"time"
)

var Cron *cron.Cron

func LoadSchedule() {
	Cron := cron.New()
	Cron.AddFunc("*/3 * * * * ?", TaskLoop, "TaskLoop")
	Cron.Start()
}

func TaskLoop() {
	conn, err := net.Dial("udp", conf.UDP_HOST)
	if err != nil {
		return
	}
	defer conn.Close()
	// 整体耗时不超过3s
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	request := "001" + conf.C.MustValue("zhimiao", "clientid")
	conn.Write([]byte(request))
	var r [20480]byte
	n, _ := conn.Read(r[0:])
	if n < 20 {
		return
	}
	resultBuf := bytes.NewBuffer(r[:n])

	var result = &dto.Hook001DTO{}
	dec := gob.NewDecoder(resultBuf)
	dec.Decode(result)
	// 判断是否含有部署数据
	if result.Deploy != "" {
		go deploy.Run(result.Deploy)
	}
}
