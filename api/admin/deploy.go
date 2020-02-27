/**
 * 部署服务
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019/6/13 16:10
 */

package admin

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common"
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/models"
	"gitee.com/zhimiao/qiansi/req"
	"gitee.com/zhimiao/qiansi/resp"
	"gitee.com/zhimiao/qiansi/udp_service"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

// @Summary 获取部署应用列表
// @Produce  json
// @Accept  json
// @Param body body req.DeployListParam true "入参集合"
// @Success 200 {array} resp.DeployVO ""
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
	if param.Id == 0 {
		member := &models.Member{
			Id: c.GetInt("UID"),
		}
		if !member.CheckDeploy() {
			resp.NewApiResult(-4, "您的部署任务创建数量已达上限").Json(c)
			return
		}
		po := &models.Deploy{}
		utils.SuperConvert(param, po)
		po.Uid = c.GetInt("UID")
		po.OpenId = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
		if models.Mysql.Save(po).RowsAffected > 0 {
			resp.NewApiResult(1, "创建成功", po).Json(c)
			return
		}
	}
	if param.Id > 0 {
		if models.Mysql.Table("deploy").Where("uid=?", c.GetInt("UID")).Save(param).RowsAffected > 0 {
			resp.NewApiResult(1, "更新成功").Json(c)
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
// @Success 200 {array} resp.DeployRunLogTabVO "返回"
// @Router /admin/DeployRunLogTab [get]
func DeployRunLogTab(c *gin.Context) {
	param := &req.DeployServerParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	d := &[]resp.DeployRunLogTabVO{}
	sql := "SELECT s.id, s.server_name,r.deploy_version FROM `server` s, `deploy_server_relation` r WHERE s.id=r.server_id and r.deploy_id=? and s.uid=?"
	models.Mysql.Raw(sql, param.DeployId, c.GetInt("UID")).Scan(d)
	resp.NewApiResult(1, "读取成功", d).Json(c)
}

// @Summary 获取当前部署应用指定服务器的运行日志
// @Produce  json
// @Accept  json
// @Param body body req.DeployRunLogParam true "入参集合"
// @Success 200 {object} resp.ApiResult ""
// @Router /admin/DeployRunLog [get]
func DeployRunLog(c *gin.Context) {
	param := &req.DeployRunLogParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	data, _ := models.InfluxDB.QueryToArray(fmt.Sprintf(
		`from(bucket: "client_log")
					|> range(start: -30d)
					|> filter(fn: (r) => r._measurement == "deploy" and r.DEPLOY_ID=="%v" and r.DEPLOY_VERSION == "%v" and r.SERVER_ID == "%v")`,
		param.DeployId,
		param.Version,
		param.ServerId,
	))
	resp.NewApiResult(1, "读取成功", data).Json(c)
}

// @Summary 获取部署日志
// @Produce  json
// @Accept  json
// @Param body body req.DeployLogParam true "入参集合"
// @Success 200 {object} resp.PageInfo ""
// @Router /admin/DeployLog [get]
func DeployLog(c *gin.Context) {
	param := &req.DeployLogParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	if param.StartTime.IsZero() {
		param.StartTime = time.Now().Add(-time.Hour * 24 * 30)
	}
	if param.EndTime.IsZero() {
		param.EndTime = time.Now()
	}
	if param.StartTime.Sub(time.Now()) > time.Hour*24*30 {
		resp.NewApiResult(-4, "日志筛选时长不可大于一个月").Json(c)
		return
	}
	fluxQuery := fmt.Sprintf(
		`from(bucket: "client_log")
					|> range(start: %s, stop: %s)
					|> filter(fn: (r) => r._measurement == "deploy" and r.SERVER_UID == "%v")`,
		param.StartTime.Format(time.RFC3339),
		param.EndTime.Format(time.RFC3339),
		c.GetInt("UID"),
	)
	if param.DeployId > 0 {
		fluxQuery += fmt.Sprintf(`|> filter(fn: (r) => r.DEPLOY_ID=="%v")`, param.DeployId)
	}
	if param.DeployVersion > 0 {
		fluxQuery += fmt.Sprintf(`|> filter(fn: (r) => r.DEPLOY_VERSION=="%v")`, param.DeployVersion)
	}
	if param.ServerId > 0 {
		fluxQuery += fmt.Sprintf(`|> filter(fn: (r) => r.SERVER_ID=="%v")`, param.ServerId)
	}
	fluxQuery += "|> group()"
	rowData, _ := models.InfluxDB.QueryToArray(fluxQuery + `|>count()`)
	var rows int
	if rowData != nil && len(rowData) > 0 {
		if v, ok := rowData[0]["_value"].(int64); ok {
			// 将 int64 转化为 int
			rows = *(*int)(unsafe.Pointer(&v))
		}
	}
	lists, _ := models.InfluxDB.QueryToArray(fluxQuery + fmt.Sprintf(`|> limit(n:%d, offset: %d)`, param.PageSize, param.Offset()))
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: rows,
		Rows:      lists,
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
	db := models.Mysql.Exec("update deploy set now_version=now_version+1 where id=?", param.DeployId)
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

// @Summary 获取部署触发链接
// @Produce  json
// @Accept  json
// @Param body body req.DeployBase true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "启动成功","data": null}"
// @Router /admin/DeployLink [GET]
func DeployLink(c *gin.Context) {
	param := &req.DeployBase{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	po := &models.Deploy{
		Id:  param.DeployId,
		Uid: c.GetInt("UID"),
	}
	if po.GetOpenId() {
		url := "/client/ApiDeployRun?Key=" + po.OpenId
		resp.NewApiResult(1, "操作成功", common.Config.Server.ApiHost+url).Json(c)
		return
	}
	resp.NewApiResult(-5, "获取失败，请重新尝试").Json(c)
}
