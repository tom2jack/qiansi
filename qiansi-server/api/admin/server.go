/**
 * 服务器模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/12 19:02
 */

package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/models"
	"qiansi/common/models/api_req"
)

// @Summary 获取服务器(客户端)列表
// @Produce  json
// @Accept  json
// @Success 200 {object} models.Server "{"code": 1,"msg": "读取成功","data": [{"CreateTime": "2019-03-02T16:10:10+08:00","DeviceId": "","Domain": "127.0.0.1","Id": 1,"ServerName": "纸喵5号机","ServerRuleId": 0,"ServerStatus": 0,"Uid": 2,"UpdateTime": "2019-05-22T17:40:18+08:00"}]}"
// @Router /admin/ServerLists [post]
func ServerLists(c *gin.Context) {
	s := &[]models.Server{}
	models.ZM_Mysql.Raw("select id, uid, server_name, server_status, server_rule_id, device_id, domain, create_time, update_time from `server` where uid=?", c.GetInt("UID")).Scan(s)
	models.NewApiResult(1, "读取成功", s).Json(c)
}

// @Summary 删除服务器
// @Produce  json
// @Accept  json
// @Param body body api_req.ServerDelParam true "入参集合"
// @Success 200 {object} models.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/ServerDel [DELETE]
func ServerDel(c *gin.Context) {

	param := &api_req.ServerDelParam{}
	if err := c.Bind(param); err != nil || param.ServerId == 0 {
		models.NewApiResult(-4, "入参解析失败").Json(c)
		return
	}
	db := models.ZM_Mysql.Delete(models.Server{}, "id=? and uid=?", param.ServerId, c.GetInt("UID"))
	if db.Error != nil || db.RowsAffected != 1 {
		models.NewApiResult(-5, "删除失败", *db).Json(c)
		return
	}
	models.NewApiResult(1, "操作成功", *db).Json(c)
}
