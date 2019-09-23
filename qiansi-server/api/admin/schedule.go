package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/utils"
	"qiansi/qiansi-server/models"
	"qiansi/qiansi-server/req"
	"qiansi/qiansi-server/resp"
	"qiansi/qiansi-server/schedule"
	"strings"
	"time"
)

// @Summary 获取计划任务列表
// @Produce  json
// @Accept  json
// @Param body body req.ScheduleListParam true "入参集合"
// @Success 200 {object} resp.PageInfo ""
// @Router /admin/ScheduleLists [get]
func ScheduleLists(c *gin.Context) {
	param := &req.ScheduleListParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	s := &models.Schedule{
		Uid:   c.GetInt("UID"),
		Title: param.Title,
	}
	lists, rows := s.List(param.Offset(), param.PageSize)
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      lists,
	}).Json(c)
}

// @Summary 创建计划任务
// @Produce  json
// @Accept  json
// @Param body body req.ScheduleCreateParam true "入参集合"
// @Success 200 {object} resp.ApiResult ""
// @Router /admin/ScheduleCreate [POST]
func ScheduleCreate(c *gin.Context) {
	param := &req.ScheduleCreateParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}

	po := &models.Schedule{}
	utils.SuperConvert(param, po)
	po.Uid = c.GetInt("UID")
	t, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	po.PrevTime = t
	po.NextTime = t
	fields := strings.Fields(po.Crontab)
	if len(fields) != 5 && len(fields) != 6 {
		resp.NewApiResult(-4, "表达式错误").Json(c)
		return
	}
	if po.Create() {
		schedule.Task.Add(po)
	}
	resp.NewApiResult(1).Json(c)
}

// @Summary 删除计划任务
// @Produce  json
// @Accept  json
// @Param body body req.ScheduleDelParam true "入参集合"
// @Success 200 {object} resp.ApiResult ""
// @Router /admin/ScheduleDel [DELETE]
func ScheduleDel(c *gin.Context) {
	param := &req.ScheduleDelParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	po := &models.Schedule{
		Id:  param.Id,
		Uid: c.GetInt("UID"),
	}
	if po.Del() {
		schedule.Task.Remove(po)
	}
	resp.NewApiResult(1).Json(c)
}

// @Summary 执行计划任务
// @Produce  json
// @Accept  json
// @Param body body req.ScheduleDoParam true "入参集合"
// @Success 200 {object} resp.ApiResult ""
// @Router /admin/ScheduleDo [get]
func ScheduleDo(c *gin.Context) {
	param := &req.ScheduleDoParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	po := &models.Schedule{
		Id:  param.Id,
		Uid: c.GetInt("UID"),
	}
	po.Get()
	if !schedule.Task.Run(po) {
		resp.NewApiResult(-5, "任务异常，无法执行").Json(c)
		return
	}
	resp.NewApiResult(1).Json(c)
}
