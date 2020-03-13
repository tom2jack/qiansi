package service

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/models"
	"time"
)

type MaxInfo struct {
	DeployNum   int
	MaxDeploy   int
	ScheduleNum int
	MaxSchedule int
}

// GetUserModuleMaxInfo 获取用户模块限额信息
func GetUserModuleMaxInfo(uid int) (info MaxInfo, err error) {
	// info = new(MaxInfo)
	member := models.Member{Id: uid}
	err = member.MaxInfo()
	if err != nil {
		return
	}
	info.MaxDeploy, info.MaxSchedule = member.MaxDeploy, member.MaxSchedule
	var num int
	deploy := models.Deploy{Uid: uid}
	num, err = deploy.Count()
	if err != nil {
		return
	}
	info.DeployNum = num
	schedule := models.Schedule{Uid: uid}
	num, err = schedule.Count()
	if err != nil {
		return
	}
	info.ScheduleNum = num
	return
}

// GetClientCPURate 获取服务器近1h的CPU利用率
func GetClientCPURate(uid, serverId int) (result []map[string]interface{}) {
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now()
	fluxQuery := fmt.Sprintf(`
		from(bucket: "client_metric")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r._measurement == "cpu" and r.QIANSI_CLIENT_ID == "%d" and r.QIANSI_CLIENT_UID == "%d")
			|> filter(fn: (r) => r._field == "usage_user" or r._field == "usage_system" or r._field == "usage_idle")
			|> filter(fn: (r) => r.cpu == "cpu-total")
			|> aggregateWindow(every: 10000ms, fn: mean, createEmpty: false)
			|> yield(name: "mean")
		`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		serverId,
		uid,
	)
	lists, _ := models.InfluxDB.QueryToArray(fluxQuery)
	return lists
}
