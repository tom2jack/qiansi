package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/models"
)

// @Summary 获取计划任务列表
// @Produce  json
// @Accept  json
// @Param body body zreq.ServerSetParam true "入参集合"
// @Success 200 {object} models.PageInfo ""
// @Router /admin/ScheduleLists [get]
func ScheduleLists(c *gin.Context) {
	s := &[]models.Server{}
	models.ZM_Mysql.Raw("select id, uid, server_name, server_status, server_rule_id, device_id, domain, create_time, update_time from `server` where uid=?", c.GetInt("UID")).Scan(s)
	models.NewApiResult(1, "读取成功", s).Json(c)
}
