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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// @Summary 获取部署应用列表
// @Produce  json
// @Accept  json
// @Param body body req.DeployListParam true "入参集合"
// @Success 200 {object} resp.ApiResult ""
// @Router /admin/DeployLists [get]
func DeployLists(c *gin.Context) {
	param := &req.DeployListParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	s := &models.Deploy{
		Uid:   c.GetInt("UID"),
		Title: param.Title,
	}
	lists, rows := s.List(param.Offset(), param.PageSize)
	vo := make([]resp.DeployVO, len(lists))
	for k, v := range lists {
		utils.SuperConvert(&v, &vo[k])
		vo[k].CreateTime = resp.JsonTimeDate(v.CreateTime)
		vo[k].UpdateTime = resp.JsonTimeDate(v.UpdateTime)
	}
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      vo,
	}).Json(c)
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
	resp.NewApiResult(1).Json(c)
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
	param := &[]req.DeployRelationParam{}
	if c.ShouldBind(param) != nil || len(*param) == 0 {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	var (
		max       = len(*param)
		serverIds = make([]int, max)
		deployIds = make([]int, max)
		result    = make([]bool, max)
	)
	for k, e := range *param {
		serverIds[k] = e.ServerId
		deployIds[k] = e.DeployId
	}
	server := &models.Server{
		Uid: c.GetInt("UID"),
	}
	if !server.BatchCheck(serverIds) {
		resp.NewApiResult(-5, "含有非法服务器绑定").Json(c)
		return
	}
	deploy := &models.Deploy{
		Uid: c.GetInt("UID"),
	}
	if !deploy.BatchCheck(deployIds) {
		resp.NewApiResult(-5, "含有非法应用绑定").Json(c)
		return
	}
	relation := &models.DeployServerRelation{}
	for k, e := range *param {
		relation.DeployId = e.DeployId
		relation.ServerId = e.ServerId
		result[k] = relation.Relation(e.Relation)
	}
	resp.NewApiResult(1, "成功", result).Json(c)
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
	server := &models.Server{
		Uid: c.GetInt("UID"),
	}
	serverList := server.ListByUser()
	relation := &models.DeployServerRelation{
		DeployId: param.DeployId,
	}
	relationList := relation.ListByDeployId()
	relationMap := make(map[int]models.DeployServerRelation)
	for _, e := range relationList {
		relationMap[e.ServerId] = e
	}
	len := len(serverList)
	list := make([]resp.DeployServerVO, len)
	for _, v := range serverList {
		d := resp.DeployServerVO{}
		len--
		utils.SuperConvert(&v, &d)
		if r, ok := relationMap[v.Id]; ok {
			utils.SuperConvert(&r, &d)
		}
		list[len] = d
	}
	resp.NewApiResult(1, "读取成功", list).Json(c)
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

// @Summary 获取部署日志
// @Produce  json
// @Accept  json
// @Param body body req.DeployLogParam true "入参集合"
// @Success 200 {array} resp.DeployLogVO "返回"
// @Router /admin/DeployLog [get]
func DeployLog(c *gin.Context) {
	param := &req.DeployLogParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	s := &models.DeployLog{
		Uid:           c.GetInt("UID"),
		DeployId:      param.DeployId,
		DeployVersion: param.DeployVersion,
		ServerId:      param.ServerId,
	}
	if param.StartTime.IsZero() {
		param.StartTime = time.Now().Add(-time.Hour * 24 * 30)
	}
	if param.EndTime.IsZero() {
		param.EndTime = time.Now()
	}
	if param.EndTime.Sub(param.StartTime) > time.Hour*24*30 {
		resp.NewApiResult(-4, "日志筛选时长不可大于一个月").Json(c)
		return
	}
	lists, rows := s.List(param.StartTime, param.EndTime, param.Offset(), param.PageSize)
	vo := make([]resp.DeployLogVO, len(lists))
	for k, v := range lists {
		utils.SuperConvert(&v, &vo[k])
		// vo[k].CreateTime = resp.JsonTimeDate(v.CreateTime)
	}
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      vo,
	}).Json(c)
}

// @Summary 启动部署
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
