package admin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/common/utils"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/service"
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
	maxInfo, _ := service.GetDashboardService().GetUserModuleMaxInfo(uid)
	vo := resp.DashboardInfoVO{
		UserMaxInfo: *maxInfo,
	}
	// vo.DeployNum
	vo.DeployRunNum, _ = models.GetDeployModels().CountDo(uid)
	// 获取服务器数量
	vo.ServerNum = models.GetServerModels().Count(uid)
	// 邀请数
	vo.InviteNum, _ = models.GetMemberModels().InviterCount(uid)
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
	uid := c.GetInt(req.UID)
	var result resp.DashboardIndexMetricVO
	result.ActiveServerNum = models.GetServerModels().GetUserActiveServerNum(uid)
	dashbardService := service.GetDashboardService()
	result.CPURate, _ = dashbardService.GetClientCPURate(uid, param.ServerId, param.StartTime, param.EndTime)
	result.MenRate, _ = dashbardService.GetClientMemRate(uid, param.ServerId, param.StartTime, param.EndTime)
	resp.NewApiResult(1, "读取成功", result).Json(c)
}
