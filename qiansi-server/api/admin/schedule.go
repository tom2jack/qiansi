package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/models"
)

// @Summary 获取计划任务列表
// @Produce  json
// @Accept  json
// @Param body body models.PageParam true "入参集合"
// @Success 200 {object} models.PageInfo ""
// @Router /admin/ScheduleLists [get]
func ScheduleLists(c *gin.Context) {
	param := &models.PageParam{}
	if err := c.ShouldBind(param); err != nil {
		models.NewApiResult(-4, "入参绑定失败"+ err.Error()).Json(c)
		return
	}
	s := &models.Schedule{
		Uid: c.GetInt("UID"),
	}
	pageInfo, err := s.List(param)
	if err != nil {
		models.NewApiResult(0, err.Error())
		return
	}
	models.NewApiResult(1, "读取成功", pageInfo).Json(c)
}
