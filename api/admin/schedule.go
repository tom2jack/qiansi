package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zhi-miao/qiansi/common/utils"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/req"
	"github.com/zhi-miao/qiansi/resp"
	"github.com/zhi-miao/qiansi/schedule"
	"github.com/zhi-miao/qiansi/service"
)

type scheduleApi struct{}

var Schedule = &scheduleApi{}

// @Summary 获取计划任务列表
// @Produce  json
// @Accept  json
// @Param body body req.ScheduleListParam true "入参集合"
// @Success 200 {object} resp.PageInfo ""
// @Router /admin/ScheduleLists [get]
func (r *scheduleApi) Lists(c *gin.Context) {
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
func (r *scheduleApi) Create(c *gin.Context) {
	param := &req.ScheduleCreateParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	info, err := service.GetUserModuleMaxInfo(c.GetInt("UID"))
	if err != nil {
		resp.NewApiResult(-4, "用户限额检测失败").Json(c)
		return
	}
	if info.MaxSchedule <= info.ScheduleNum {
		resp.NewApiResult(-4, "您的调度任务创建数量已达上限").Json(c)
		return
	}
	po := &models.Schedule{}
	utils.SuperConvert(param, po)
	po.Uid = c.GetInt("UID")
	err = utils.PanicToError(func() {
		po.NextTime = schedule.Task.NextTime(po.Crontab)
	})
	if err != nil {
		resp.NewApiResult(-4, "表达式错误").Json(c)
		return
	}
	if !po.Create() {
		resp.NewApiResult(-5, "任务档案建立失败").Json(c)
		return
	}
	err = schedule.Task.Add(*po)
	if err != nil {
		resp.NewApiResult(-5, "添加任务到调度器失败").Json(c)
		return
	}
	resp.NewApiResult(1).Json(c)
}

// @Summary 删除计划任务
// @Produce  json
// @Accept  json
// @Param body body req.ScheduleDelParam true "入参集合"
// @Success 200 {object} resp.ApiResult ""
// @Router /admin/ScheduleDel [DELETE]
func (r *scheduleApi) Del(c *gin.Context) {
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
func (r *scheduleApi) Do(c *gin.Context) {
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
