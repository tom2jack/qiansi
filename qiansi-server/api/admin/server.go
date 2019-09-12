/**
 * 服务器模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/12 19:02
 */

package admin

import (
	"qiansi/common/utils"
	"qiansi/qiansi-server/models"
	"qiansi/qiansi-server/req"
	"qiansi/qiansi-server/resp"

	"github.com/gin-gonic/gin"
)

// @Summary 获取服务器(客户端)列表
// @Produce  json
// @Accept  json
// @Param body body req.ServerListParam true "入参集合"
// @Success 200 {array} resp.ServerVO ""
// @Router /admin/ServerLists [get]
func ServerLists(c *gin.Context) {
	param := &req.ServerListParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	s := &models.Server{
		Uid:        c.GetInt("UID"),
		ServerName: param.ServerName,
	}
	lists, rows := s.List(param.Offset(), param.PageSize)
	vo := make([]resp.ServerVO, len(lists))
	for k, v := range lists {
		utils.SuperConvert(&v, &vo[k])
	}
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      vo,
	}).Json(c)
}

// @Summary 设置服务器信息
// @Produce  json
// @Accept  json
// @Param body body req.ServerSetParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/ServerSet [POST]
func ServerSet(c *gin.Context) {
	param := &req.ServerSetParam{}
	if c.ShouldBind(param) != nil {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	po := &models.Server{}
	utils.SuperConvert(param, po)
	if models.Mysql.Table("server").Where("id=? and uid=?", po.Id, c.GetInt("UID")).Updates(po).RowsAffected > 0 {
		resp.NewApiResult(1, "更新成功", po).Json(c)
		return
	}
	resp.NewApiResult(0, "系统错误").Json(c)
}

// @Summary 删除服务器
// @Produce  json
// @Accept  json
// @Param body body req.ServerDelParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/ServerDel [DELETE]
func ServerDel(c *gin.Context) {
	param := &req.ServerDelParam{}
	if err := c.Bind(param); err != nil || param.ServerId == 0 {
		resp.NewApiResult(-4, "入参解析失败").Json(c)
		return
	}
	db := models.Mysql.Delete(models.Server{}, "id=? and uid=?", param.ServerId, c.GetInt("UID"))
	if db.Error != nil || db.RowsAffected != 1 {
		resp.NewApiResult(-5, "删除失败", db).Json(c)
		return
	}
	// 删除关联
	models.Mysql.Delete(models.DeployServerRelation{}, "server_id=?", param.ServerId)
	resp.NewApiResult(1, "操作成功", db).Json(c)
}
