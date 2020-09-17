/**
 * 服务器模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/12 19:02
 */

package admin

import (
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"

	"github.com/gin-gonic/gin"
)

type serverApi struct{}

var Server = &serverApi{}

// @Summary 获取服务器(客户端)列表
// @Produce  json
// @Accept  json
// @Param body body req.ServerListParam true "入参集合"
// @Success 200 {array} resp.ServerVO ""
// @Router /admin/ServerLists [get]
func (r *serverApi) Lists(c *gin.Context) {
	param := &req.ServerListParam{}
	if err := c.ShouldBind(param); err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
	uid := c.GetInt(req.UID)
	lists, rows := models.GetServerModels().List(uid, param)
	vo := make([]resp.ServerVO, len(lists))
	for k, v := range lists {
		gutils.SuperConvert(&v, &vo[k])
	}
	c.JSON(resp.ApiSuccess(resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      vo,
	}))
}

// @Summary 设置服务器信息
// @Produce  json
// @Accept  json
// @Param body body req.ServerSetParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/ServerSet [POST]
func (r *serverApi) Set(c *gin.Context) {
	param := &req.ServerSetParam{}
	if c.ShouldBind(param) != nil {
		c.JSON(resp.ApiError("入参绑定失败"))
		return
	}
	po := &models.Server{}
	gutils.SuperConvert(param, po)
	// TODO: 移动至models
	err := models.Mysql.Table("server").Where("id=? and uid=?", po.ID, c.GetInt("UID")).Updates(po).Error
	if err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
}

// @Summary 删除服务器
// @Produce  json
// @Accept  json
// @Param body body req.ServerDelParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/ServerDel [DELETE]
func (r *serverApi) Del(c *gin.Context) {
	param := &req.ServerDelParam{}
	if err := c.Bind(param); err != nil || param.ServerId == 0 {
		c.JSON(resp.ApiError("入参解析失败"))
		return
	}
	// TODO: 移动至models
	db := models.Mysql.Delete(models.Server{}, "id=? and uid=?", param.ServerId, c.GetInt("UID"))
	if db.Error != nil || db.RowsAffected != 1 {
		c.JSON(resp.ApiError("删除失败"))
		return
	}
	// 删除关联
	models.Mysql.Delete(models.DeployServerRelation{}, "server_id=?", param.ServerId)
}
