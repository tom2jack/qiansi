package admin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhi-miao/qiansi/common/utils"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/notifyevent"
	"github.com/zhi-miao/qiansi/req"
	"github.com/zhi-miao/qiansi/resp"
	"github.com/zhi-miao/qiansi/service"
	"time"
)

type dashboardApi struct{}

var Dashboard = &dashboardApi{}

// @Summary 获取基本信息
// @Produce  json
// @Accept  json
// @Success 200 {array} resp.DashboardInfoVO ""
// @Router /admin/dashboard/info [get]
func (r *dashboardApi) Info(c *gin.Context) {
	uid := c.GetInt("UID")
	// 获取限额信息
	maxInfo, _ := service.GetUserModuleMaxInfo(uid)
	vo := resp.DashboardInfoVO{
		DeployNum:   maxInfo.DeployNum,
		MaxDeploy:   maxInfo.MaxDeploy,
		ScheduleNum: maxInfo.ScheduleNum,
		MaxSchedule: maxInfo.MaxSchedule,
	}
	// 部署次數
	deploy := &models.Deploy{UId: uid}
	// vo.DeployNum
	vo.DeployRunNum, _ = deploy.CountDo()
	// 获取服务器数量
	vo.ServerNum = models.GetServerModels().Count(uid)
	// 邀请数
	member := &models.Member{InviterUid: uid}
	vo.InviteNum, _ = member.InviterCount()
	resp.NewApiResult(1, "读取成功", vo).Json(c)
}

// @Summary 概览大盘
// @Produce  json
// @Accept  json
// @Param body body req.DashboardIndexMetricParam true "入参集合"
// @Success 200 {array} resp.DashboardIndexMetricVO ""
// @Router /admin/dashboard/IndexMetric [get]
func (r *dashboardApi) IndexMetric(c *gin.Context) {
	param := &req.DashboardIndexMetricParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	if param.StartTime.IsZero() {
		param.EndTime = time.Now()
	}
	if param.EndTime.IsZero() {
		param.StartTime = param.EndTime.Add(-1 * time.Hour)
	}
	if param.EndTime.Sub(param.StartTime) > 2*time.Hour {
		resp.NewApiResult(-4, "检索时间不可大于2H").Json(c)
		return
	}
	uid := c.GetInt("UID")
	serIdsCacheId := fmt.Sprintf("QIANSI:dashboard:user-server-ids:%d", uid)
	var serIds []int
	s, err := models.Redis.Get(serIdsCacheId)
	if err == nil && s != "" {
		json.Unmarshal([]byte(s), &serIds)
	} else {
		serIds = models.GetServerModels().UserServerIds(uid)
		if serIdsJson, err := json.Marshal(serIds); err == nil {
			models.Redis.Set(serIdsCacheId, string(serIdsJson), 5*60)
		}
	}
	vo := resp.DashboardIndexMetricVO{
		ActiveServerNum: notifyevent.Hook001.GetActiveServerNum(serIds...),
		CPURate:         service.GetClientCPURate(uid, param.ServerId, param.StartTime, param.EndTime),
		MenRate:         service.GetClientMemRate(uid, param.ServerId, param.StartTime, param.EndTime),
	}
	resp.NewApiResult(1, "读取成功", vo).Json(c)
}
