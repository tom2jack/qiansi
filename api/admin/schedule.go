package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"
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
		c.JSON(resp.ApiError(err))
		return
	}
	uid := c.GetInt(req.UID)
	lists, rows := models.GetScheduleModels().List(uid, param)
	c.JSON(resp.ApiSuccess(resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      lists,
	}))

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
		c.JSON(resp.ApiError(err))
		return
	}
	uid := c.GetInt(req.UID)
	err := service.GetScheduleService().Create(uid, param)
	if err != nil {
		c.JSON(resp.ApiError(err))
	}
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
		c.JSON(resp.ApiError(err))
		return
	}
	uid := c.GetInt(req.UID)
	err := models.GetScheduleModels().Del(param.Id, uid)
	if err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
	schedule.Task.Remove(param.Id)
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
		c.JSON(resp.ApiError(err))
		return
	}
	uid := c.GetInt(req.UID)
	po, err := models.GetScheduleModels().Get(param.Id, uid)
	if err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
	if !schedule.Task.Run(&po) {
		c.JSON(resp.ApiError("任务异常，无法执行"))
		return
	}
}
