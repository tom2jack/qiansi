package admin

import (
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/qiansi-server/models"
	"gitee.com/zhimiao/qiansi/qiansi-server/req"
	"gitee.com/zhimiao/qiansi/qiansi-server/resp"
	"gitee.com/zhimiao/qiansi/qiansi-server/schedule"
	"github.com/gin-gonic/gin"
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
	member := &models.Member{
		Id: c.GetInt("UID"),
	}
	if !member.CheckSchedule() {
		resp.NewApiResult(-4, "您的调度任务创建数量已达上限").Json(c)
		return
	}
	po := &models.Schedule{}
	utils.SuperConvert(param, po)
	po.Uid = c.GetInt("UID")
	err := utils.PanicToError(func() {
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
