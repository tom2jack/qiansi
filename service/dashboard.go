package service

import (
	"fmt"
	"time"

	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"
)

type dashboardService struct{}

func GetDashboardService() *dashboardService {
	return &dashboardService{}
}

// GetUserModuleMaxInfo 获取用户模块限额信息
func (s *dashboardService) GetUserModuleMaxInfo(uid int) (*resp.UserMaxInfo, error) {
	var info resp.UserMaxInfo
	user, err := models.GetMemberModels().UsereInfo(uid)
	if err != nil {
		return nil, err
	}
	info.MaxDeploy, info.MaxSchedule = user.MaxDeploy, user.MaxSchedule
	info.DeployNum, err = models.GetDeployModels().Count(uid)
	if err != nil {
		return nil, err
	}
	info.ScheduleNum, err = models.GetScheduleModels().Count(uid)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// GetClientCPURate 获取服务器近1h的CPU利用率
func (s *dashboardService) GetClientCPURate(uid, serverId int, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	fluxQuery := fmt.Sprintf(`
		from(bucket: "client_metric")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["SERVER_ID"] == "%d" and r["SERVER_UID"] == "%d")
			|> filter(fn: (r) => r["_measurement"] == "cpu")
			|> filter(fn: (r) => r["_field"] == "usage_user")
			|> filter(fn: (r) => r["cpu"] == "cpu-total")
		`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		serverId,
		uid,
	)
	return models.InfluxDB.QueryToArray(fluxQuery)
}

// GetClientMemRate 获取服务器近1h的内存利用率
func (s *dashboardService) GetClientMemRate(uid, serverId int, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	fluxQuery := fmt.Sprintf(`
		from(bucket: "client_metric")
			|> range(start: %s, stop: %s)
			|> filter(fn: (r) => r["SERVER_ID"] == "%d" and r["SERVER_UID"] == "%d")
			|> filter(fn: (r) => r["_measurement"] == "mem")
			|> filter(fn: (r) => r["_field"] == "used_percent")
		`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		serverId,
		uid,
	)
	return models.InfluxDB.QueryToArray(fluxQuery)
}
