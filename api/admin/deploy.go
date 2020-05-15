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
	"gitee.com/zhimiao/qiansi/notifyevent"
	"gitee.com/zhimiao/qiansi/req"
	"gitee.com/zhimiao/qiansi/resp"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

type deployApi struct{}

var Deploy = &deployApi{}

// @Summary 获取部署应用列表
// @Produce  json
// @Accept  json
// @Param body body req.DeployListParam true "入参集合"
// @Success 200 {array} models.Deploy ""
// @Router /admin/DeployLists [get]
func (r *deployApi) Lists(c *gin.Context) {
	param := &req.DeployListParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	list, totalRows := models.DeployList(c.GetInt("UID"), param)
	resp.NewApiResult(1, "读取成功", resp.PageInfo{
		Page:      param.Page,
		PageSize:  param.PageSize,
		TotalSize: totalRows,
		Rows:      list,
	}).Json(c)
}

// @Summary 创建部署应用
// @Produce  json
// @Accept  json
// @Param body body req.DeploySetParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/Deploy [POST]
func (r *deployApi) Create(c *gin.Context) {
	param := &req.DeploySetParam{}
	if c.ShouldBind(param) != nil {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	err := models.CreateDeploy(c.GetInt("UID"), param)
	resp.NewApiResult(err).Json(c)
}

// @Summary 更新部署应用
// @Produce  json
// @Accept  json
// @Param body body req.DeploySetParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/Deploy [put]
func (r *deployApi) Update(c *gin.Context) {
	param := &req.DeploySetParam{}
	if c.ShouldBind(param) != nil {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	err := models.UpdateDeploy(c.GetInt("UID"), param)
	resp.NewApiResult(err).Json(c)
}

// @Summary 删除部署应用
// @Produce  json
// @Accept  json
// @Param body body req.DeployParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功","data": null}"
// @Router /admin/Deploy [DELETE]
func (r *deployApi) Del(c *gin.Context) {
	param := &req.DeployParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	err := models.DelDeploy(c.GetInt("UID"), param.DeployId)
	resp.NewApiResult(err).Json(c)
}

// @Summary 获取当前部署应用的服务器列表
// @Produce  json
// @Accept  json
// @Param body body req.DeployParam true "入参集合"
// @Success 200 {array} models.DeployServerRelation "返回"
// @Router /admin/DeployServer [post]
func (r *deployApi) DeployServer(c *gin.Context) {
	param := &req.DeployParam{}
	if c.ShouldBind(param) != nil || !(param.DeployId > 0) {
		resp.NewApiResult(-4, "入参绑定失败").Json(c)
		return
	}
	list, err := models.DeployServerList(c.GetInt("UID"), param.DeployId)
	if err != nil {
		resp.NewApiResult(err).Json(c)
		return
	}
	resp.NewApiResult(1, "读取成功", list).Json(c)
}

// @Summary 获取当前部署应用绑定的服务器列表，用于渲染运行日志选项卡
// @Produce  json
// @Accept  json
// @Param body body req.DeployServerParam true "入参集合"
// @Success 200 {array} resp.DeployRunLogTabVO "返回"
// @Router /admin/DeployRunLogTab [get]
func (r *deployApi) RunLogTab(c *gin.Context) {
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
func (r *deployApi) RunLog(c *gin.Context) {
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
func (r *deployApi) Log(c *gin.Context) {
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
func (r *deployApi) Do(c *gin.Context) {
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
		notifyevent.Hook001.AddDeploy(v.Id)
	}
	resp.NewApiResult(1, "启动成功", server).Json(c)
}

// @Summary 获取部署触发链接
// @Produce  json
// @Accept  json
// @Param body body req.DeployParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "启动成功","data": null}"
// @Router /admin/DeployLink [GET]
func (r *deployApi) Link(c *gin.Context) {
	param := &req.DeployParam{}
	if err := c.ShouldBind(param); err != nil {
		resp.NewApiResult(-4, utils.Validator(err)).Json(c)
		return
	}
	po := &models.Deploy{
		ID:  param.DeployId,
		UId: c.GetInt("UID"),
	}
	if po.GetOpenId() {
		url := "/client/ApiDeployRun?Key=" + po.OpenID
		resp.NewApiResult(1, "操作成功", common.Config.Server.ApiHost+url).Json(c)
		return
	}
	resp.NewApiResult(-5, "获取失败，请重新尝试").Json(c)
}
