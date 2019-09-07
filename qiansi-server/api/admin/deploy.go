/**
 * 部署服务
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/13 16:10
 */

package admin

import (
	"qiansi/common/utils"
	"qiansi/qiansi-server/models"
	"qiansi/qiansi-server/req"
	"qiansi/qiansi-server/resp"
	"qiansi/qiansi-server/udp_service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// @Summary 获取部署应用列表
// @Produce  json
// @Accept  json
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "读取成功","data": [{"AfterCommand": "324545","BeforeCommand": "1232132132","Branch": "123213","CreateTime": "2019-02-28T10:24:41+08:00","DeployType": 1,"Id": 491,"LocalPath": "123213","NowVersion": 0,"RemoteUrl": "123213","Title": "491-一号机器的修改241","Uid": 2,"UpdateTime": "2019-02-28T10:25:17+08:00"}]}"
// @Router /admin/DeployLists [get]
func DeployLists(c *gin.Context) {
	vo := []resp.DeployVO{}
	models.Mysql.Raw("select * from deploy where uid=? order by id desc", c.GetInt("UID")).Scan(&vo)
	resp.NewApiResult(1, "读取成功", vo).Json(c)
}

// @Summary 设置部署应用
// @Produce  json
// @Accept  json
// @Param body body req.DeploySetParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/DeploySet [POST]
func DeploySet(c *gin.Context) {
	param := &req.DeploySetParam{}
	if c.ShouldBind(param) != nil {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	po := &models.Deploy{}
	utils.SuperConvert(param, po)

	if po.Id == 0 {
		po.Uid = c.GetInt("UID")
		if models.Mysql.Save(po).RowsAffected > 0 {
			resp.NewApiResult(1, "创建成功", po).Json(c)
			return
		}
	}
	if po.Id > 0 {
		if models.Mysql.Table("deploy").Where("id=? and uid=?", po.Id, c.GetInt("UID")).Updates(po).RowsAffected > 0 {
			resp.NewApiResult(1, "更新成功", po).Json(c)
			return
		}
	}
	resp.NewApiResult(0, "系统错误").Json(c)
}

// @Summary 删除部署应用
// @Produce  json
// @Accept  json
// @Param body body req.DeployDelParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/DeployDel [DELETE]
func DeployDel(c *gin.Context) {
	param := &req.DeployDelParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	db := models.Mysql.Delete(&models.Deploy{}, "id=? and uid=?", param.DeployId, c.GetInt("UID"))
	if db.Error != nil || db.RowsAffected != 1 {
		resp.NewApiResult(-5, "删除失败", db).Json(c)
		return
	}
	// 删除关联
	models.Mysql.Delete(&models.DeployServerRelation{}, "deploy_id=?", param.DeployId)
	resp.NewApiResult(1, "操作成功", db).Json(c)
}

// @Summary 部署应用关联服务器
// @Produce  json
// @Accept  json
// @Param body body req.DeployRelationParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "关联成功","data": null}"
// @Router /admin/DeployRelationServer [POST]
func DeployRelationServer(c *gin.Context) {
	param := &req.DeployRelationParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) || !(param.ServerId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	var (
		num int
		db  *gorm.DB
	)
	db = models.Mysql.Table("server").Where("id=? and uid=?", param.ServerId, c.GetInt("UID")).Count(&num)
	if db.Error != nil || num == 0 {
		resp.NewApiResult(-5, "服务器不存在").Json(c)
		return
	}
	db = models.Mysql.Table("deploy").Where("id=? and uid=?", param.DeployId, c.GetInt("UID")).Count(&num)
	if db.Error != nil || num == 0 {
		resp.NewApiResult(-5, "部署服务不存在").Json(c)
		return
	}
	if param.Relation {
		db = models.Mysql.Save(&models.DeployServerRelation{
			ServerId: param.ServerId,
			DeployId: param.DeployId,
		})
	} else {
		db = models.Mysql.Delete(&models.DeployServerRelation{}, "server_id=? and deploy_id=?", param.ServerId, param.DeployId)
	}
	if db.Error != nil || db.RowsAffected != 1 {
		resp.NewApiResult(-5, "失败").Json(c)
		return
	}
	resp.NewApiResult(1, "成功").Json(c)
}

// @Summary 获取当前部署应用的服务器列表
// @Produce  json
// @Accept  json
// @Param body body req.DeployServerParam true "入参集合"
// @Success 200 {object} models.Server "返回"
// @Router /admin/DeployServer [post]
func DeployServer(c *gin.Context) {
	param := &req.DeployServerParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	d := &[]resp.DeployServerVO{}
	sql := "SELECT s.*,r.deploy_version FROM `server` s, `deploy_server_relation` r WHERE s.id=r.server_id and r.deploy_id=? and s.uid=?"
	models.Mysql.Raw(sql, param.DeployId, c.GetInt("UID")).Scan(d)
	resp.NewApiResult(1, "读取成功", d).Json(c)
}

// @Summary 获取当前部署应用绑定的服务器列表，用于渲染运行日志选项卡
// @Produce  json
// @Accept  json
// @Param body body req.DeployServerParam true "入参集合"
// @Success 200 {object} models.Server "返回"
// @Router /admin/DeployRunLogTab [get]
func DeployRunLogTab(c *gin.Context) {
	param := &req.DeployServerParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	d := &[]resp.DeployServerVO{}
	sql := "SELECT s.id, s.server_name,r.deploy_version FROM `server` s, `deploy_server_relation` r WHERE s.id=r.server_id and r.deploy_id=? and s.uid=?"
	models.Mysql.Raw(sql, param.DeployId, c.GetInt("UID")).Scan(d)
	resp.NewApiResult(1, "读取成功", d).Json(c)
}

// @Summary 获取当前部署应用指定服务器的运行日志
// @Produce  json
// @Accept  json
// @Param body body req.DeployRunLogParam true "入参集合"
// @Success 200 {object} models.DeployLog "返回"
// @Router /admin/DeployRunLog [get]
func DeployRunLog(c *gin.Context) {
	param := &req.DeployRunLogParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	deployLogs := &[]models.DeployLog{}
	models.Mysql.Where("uid=? and deploy_id=? and server_id=? and deploy_version=?", c.GetInt("UID"), param.DeployId, param.ServerId, param.Version).Find(deployLogs)
	resp.NewApiResult(1, "读取成功", deployLogs).Json(c)
}

// @Summary 启动部署 TODO: 后期关闭此接口的开放特性，新增外部接口，通过不可枚举key作为部署参数
// @Produce  json
// @Accept  json
// @Param body body req.DeployDoParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "启动成功","data": null}"
// @Router /admin/DeployDo [GET]
func DeployDo(c *gin.Context) {
	param := &req.DeployDoParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	var (
		db *gorm.DB
	)
	db = models.Mysql.Exec("update deploy set now_version=now_version+1 where id=?", param.DeployId)
	if db.Error != nil || db.RowsAffected != 1 {
		resp.NewApiResult(-5, "部署服务不存在").Json(c)
		return
	}
	server := &[]models.Server{}
	models.Mysql.Select("id").Where("id in (select server_id from deploy_server_relation where deploy_id=?)", param.DeployId).Find(server)
	for _, v := range *server {
		udp_service.Hook001.Deploy.SET(strconv.Itoa(v.Id), "1")
	}
	resp.NewApiResult(1, "启动成功", server).Json(c)
}
