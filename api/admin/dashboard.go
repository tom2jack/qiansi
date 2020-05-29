package admin

import (
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/models"
	"gitee.com/zhimiao/qiansi/notifyevent"
	"gitee.com/zhimiao/qiansi/req"
	"gitee.com/zhimiao/qiansi/resp"
	"gitee.com/zhimiao/qiansi/service"
	"github.com/gin-gonic/gin"
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
	server := &models.Server{Uid: uid}
	vo.ServerNum, _ = server.Count()
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
	param.EndTime = time.Now()
	param.StartTime = param.EndTime.Add(-1 * time.Minute)
	uid := c.GetInt("UID")
	s := &models.Server{Uid: uid}
	serIds := s.UserServerIds()
	vo := resp.DashboardIndexMetricVO{
		ActiveServerNum: notifyevent.Hook001.GetActiveServerNum(serIds...),
		CPURate:         service.GetClientCPURate(uid, param.ServerId, param.StartTime, param.EndTime),
		MenRate:         service.GetClientMemRate(uid, param.ServerId, param.StartTime, param.EndTime),
	}
	resp.NewApiResult(1, "读取成功", vo).Json(c)
}
