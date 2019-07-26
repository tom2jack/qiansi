package schedule

import (
	"github.com/jakecoffman/cron"
	"net"
	"qiansi/common/conf"
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
	choose := string(r[:3])
	data := r[3:n]
	switch choose {
	case "001":
		deploy.Run(data)
	}
}
