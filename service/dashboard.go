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
	member := &models.Member{Id: uid}
	err = member.MaxInfo()
	if err != nil {
		return
	}
	info.MaxDeploy, info.MaxSchedule = member.MaxDeploy, member.MaxSchedule
	var num int
	deploy := &models.Deploy{UId: uid}
	num, err = deploy.Count()
	if err != nil {
		return
	}
	info.DeployNum = num
	schedule := &models.Schedule{Uid: uid}
	num, err = schedule.Count()
	if err != nil {
		return
	}
	info.ScheduleNum = num
	return
}

// GetClientCPURate 获取服务器近1h的CPU利用率
func GetClientCPURate(uid, serverId int, startTime, endTime time.Time) (result []map[string]interface{}) {
	fluxQuery := fmt.Sprintf(`
		from(bucket: "client_metric")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r._measurement == "cpu" and r.SERVER_ID == "%d" and r.SERVER_UID == "%d")
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

// GetClientMemRate 获取服务器近1h的内存利用率
func GetClientMemRate(uid, serverId int, startTime, endTime time.Time) (result []map[string]interface{}) {
	fluxQuery := fmt.Sprintf(`
		from(bucket: "client_metric")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r._measurement == "mem" and r.SERVER_ID == "%d" and r.SERVER_UID == "%d")
			|> filter(fn: (r) => r._field == "used_percent")
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
