package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/utils"
	"qiansi/qiansi-server/models"
	"qiansi/qiansi-server/req"
	"qiansi/qiansi-server/resp"
)

// @Summary 获取计划任务列表
// @Produce  json
// @Accept  json
// @Param body body req.PageParam true "入参集合"
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
		PageSize:  len(lists),
		TotalSize: rows,
		Rows:      lists,
	}).Json(c)
}
