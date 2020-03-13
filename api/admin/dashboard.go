package admin

import (
	"gitee.com/zhimiao/qiansi/models"
	"gitee.com/zhimiao/qiansi/notifyevent"
	"gitee.com/zhimiao/qiansi/resp"
	"gitee.com/zhimiao/qiansi/service"
	"github.com/gin-gonic/gin"
	"strconv"
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
	deploy := &models.Deploy{Uid: uid}
	vo.DeployRunNum, _ = deploy.CountDo()
	// 获取服务器数量
	server := &models.Server{Uid: uid}
	vo.ServerNum, _ = server.Count()
	// 邀请数
	member := models.Member{Id: uid}
	vo.InviteNum, _ = member.InviterCount()
	resp.NewApiResult(1, "读取成功", vo).Json(c)
}

// @Summary 概览大盘
// @Produce  json
// @Accept  json
// @Param body body req.DeployListParam true "入参集合"
// @Success 200 {array} resp.DeployVO ""
// @Router /admin/dashboard/IndexMetric [get]
func (r *dashboardApi) IndexMetric(c *gin.Context) {
	serverId, err := strconv.Atoi(c.Query("serverId"))
	if err != nil || serverId == 0 {
		resp.NewApiResult(-4, "服务器ID错误").Json(c)
		return
	}
	uid := c.GetInt("UID")
	s := &models.Server{Uid: uid}
	serIds := s.UserServerIds()
	activeServerNum := notifyevent.Hook001.GetActiveServerNum(serIds...)
	print(activeServerNum)

	cpuRate := service.GetClientCPURate(uid, serverId)
	resp.NewApiResult(1, "读取成功", cpuRate).Json(c)
}
