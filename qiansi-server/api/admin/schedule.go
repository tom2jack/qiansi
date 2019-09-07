package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
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
	param := &req.PageParam{}
	if err := c.ShouldBind(param); err != nil {
		if v8, ok := err.(validator.ValidationErrors); ok {
			for _, v := range v8 {
				resp.NewApiResult(-4, fmt.Sprintf("%s参数%s规则校验失败", v.Field, v.ActualTag)).Json(c)
				return
			}
		}
		resp.NewApiResult(-4, "入参绑定失败"+ err.Error()).Json(c)
		return
	}
	s := &models.Schedule{
		Uid: c.GetInt("UID"),
	}
	lists, rows := s.List(param.Offset(), param.PageSize)
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      lists,
	}).Json(c)
}
