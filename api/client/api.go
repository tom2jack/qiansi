package client

import (
	"encoding/json"
	"fmt"
	"gitee.com/zhimiao/qiansi/dto"
	"gitee.com/zhimiao/qiansi/models"
	"gitee.com/zhimiao/qiansi/resp"
	"gitee.com/zhimiao/qiansi/udp_service"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb-client-go"
	"github.com/lifei6671/gorand"
	"strconv"
	"time"
)

// @Summary 服务器注册
// @Produce  json
// @Accept  json
// @Param uid query string true "用户UID"
// @Param device query string true "客户端设备号"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /clinet/ApiRegServer [get]
func ApiRegServer(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("uid"))
	if !(uid > 0) {
		resp.NewApiResult(-4, "用户UID非法").Json(c)
		return
	}
	device := c.Query("device")
	if len(device) != 36 {
		resp.NewApiResult(-4, "客户端唯一标识号非法").Json(c)
		return
	}
	var row int
	models.Mysql.Table("member").Where("id = ?", uid).Count(&row)
	if row == 0 {
		resp.NewApiResult(-5, "用户不存在").Json(c)
		return
	}
	models.Mysql.Table("server").Where("device_id=?", device).Count(&row)
	if row > 0 {
		resp.NewApiResult(-5, "设备已存在，请勿重复注册").Json(c)
		return
	}
	api_secret := string(gorand.KRand(16, gorand.KC_RAND_KIND_ALL))
	server := &models.Server{
		Uid:          uid,
		ApiSecret:    api_secret,
		DeviceId:     device,
		ServerStatus: 1,
		Domain:       c.ClientIP(),
	}
	models.Mysql.Create(server)
	resp.NewApiResult(1, "成功", server).Json(c)
}

// @Summary 获取服务器部署任务清单
// @Produce  json
// @Accept  json
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "读取成功","data": [deploy]}"
// @Router /clinet/ApiGetDeployTask [GET]
func ApiGetDeployTask(c *gin.Context) {
	server_id := c.GetInt("SERVER-ID")
	defer udp_service.Hook001.DelDeploy(server_id)
	deploy := &[]models.Deploy{}
	models.Mysql.Raw("SELECT d.* FROM `deploy` d LEFT JOIN `deploy_server_relation` r ON d.id=r.deploy_id WHERE r.server_id=? and d.now_version > r.deploy_version", server_id).Scan(deploy)
	resp.NewApiResult(1, "读取成功", deploy).Encypt(c)
}

// @Summary 获取Telegraf监控配置
// @Produce  json
// @Accept  json
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "读取成功","data": [deploy]}"
// @Router /clinet/ApiGetTelegrafConfig [GET]
func ApiGetTelegrafConfig(c *gin.Context) {
	server_id := c.GetInt("SERVER-ID")
	uid := c.GetInt("SERVER-UID")
	sysConfig := &models.SysConfig{}
	sysConfig.Get("telegraf_config")
	defaultConfig := sysConfig.Data
	telegraf := &models.Telegraf{
		ServerID: server_id,
		UId:      uid,
	}
	telegraf.Get()
	// TODO: 配置合并
	selfConfig := telegraf.TomlConfig
	if selfConfig == "" {
		selfConfig = defaultConfig
	}
	// 判断当前客户端是否开启监控
	var isOpen = true
	if telegraf.IsOpen == 2 {
		isOpen = false
	}
	resuslt := map[string]interface{}{
		"toml_config": selfConfig,
		"is_open":     isOpen,
	}
	// 清理更新消息通知
	udp_service.Hook001.DelTelegraf(server_id)
	resp.NewApiResult(1, "读取成功", resuslt).Encypt(c)
}

// @Summary 客户端监控指标推送
// @Produce  json
// @Accept  json
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "操作成功"}"
// @Router /clinet/ApiClientMetric [post]
func ApiClientMetric(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.Status(403)
		return
	}
	type dataModel struct {
		Metrics []dto.ClientMetricDTO `json:"metrics"`
	}
	rawData := dataModel{}
	err = json.Unmarshal(raw, &rawData)
	if err != nil {
		fmt.Print(err.Error())
		c.Status(400)
		return
	}
	if len(rawData.Metrics) == 0 {
		resp.NewApiResult(1).Json(c)
		return
	}
	mds := make([]influxdb.Metric, len(rawData.Metrics))
	for k, v := range rawData.Metrics {
		mds[k] = influxdb.NewRowMetric(v.Fields, v.Name, v.Tags, time.Unix(v.Timestamp, 0))
	}
	models.InfluxDB.Write("client_metric", mds...)
	resp.NewApiResult(1).Json(c)
}

// @Summary 客户端日志推送
// @Produce  json
// @Accept  multipart/form-data
// @Param deployId formData string true "部署应用ID"
// @Param version formData string true "部署版本号"
// @Param content formData string true "日志文本内容"
// @Success 200 {object} resp.ApiResult ""
// @Router /clinet/ApiDeployLog [post]
func ApiDeployLog(c *gin.Context) {
	raw, err := c.GetRawData()
	if err != nil {
		c.Status(403)
		return
	}
	rawData := map[string]interface{}{}
	err = json.Unmarshal(raw, &rawData)
	if err != nil {
		c.Status(400)
		return
	}
	mFields := map[string]interface{}{}
	mName := ""
	mTags := map[string]string{
		"SERVER_ID":     strconv.Itoa(c.GetInt("SERVER-ID")),
		"SERVER_UID":    strconv.Itoa(c.GetInt("SERVER-UID")),
		"SERVER_DEVICE": c.GetString("SERVER-DEVICE"),
	}
	for k, v := range rawData {
		vData := fmt.Sprintf("%v", v)
		switch k {
		case "__MODEL__":
			mName = vData
		case "__MESSAGE__":
			mFields["Message"] = vData
		case "__LOG_LEVEL__":
			mTags["LOG_LEVEL"] = vData
		default:
			mTags[k] = vData
		}
	}
	metric := influxdb.NewRowMetric(mFields, mName, mTags, time.Now())
	models.InfluxDB.Write("client_log", metric)
	resp.NewApiResult(1).Json(c)
}

// @Summary 客户端部署成功回调
// @Produce  json
// @Accept  json
// @Param version query string true "版本号"
// @Param deploy_id query string true "部署应用ID"
// @Success 200 {object} resp.ApiResult ""
// @Router /clinet/ApiDeployNotify [get]
func ApiDeployNotify(c *gin.Context) {
	serverId := c.GetInt("SERVER-ID")
	version, _ := strconv.Atoi(c.Query("version"))
	deployId, _ := strconv.Atoi(c.Query("deployId"))
	uid := c.GetInt("SERVER-UID")
	if version > 0 {
		var nowVersion int
		row := models.Mysql.Raw("select now_version from deploy where id=? and uid=?", deployId, uid).Row()
		row.Scan(&nowVersion)
		if nowVersion >= version {
			if models.Mysql.Exec("update deploy_server_relation set deploy_version=? where deploy_id=? and server_id=?", version, deployId, serverId).RowsAffected > 0 {
				resp.NewApiResult(1).Json(c)
				return
			}
		}
	}
	resp.NewApiResult(1, "没改成功").Json(c)
}

// @Summary 启动部署
// @Produce  json
// @Accept  json
// @Param key query string true "入参集合"
// @Success 200 {object} string "操作结果"
// @Router /client/ApiDeployRun [GET]
func ApiDeployRun(c *gin.Context) {
	openId := c.Query("Key")
	if len(openId) != 32 {
		c.String(404, "服务不存在")
		return
	}
	deployPO := &models.Deploy{
		OpenId: openId,
	}
	if !deployPO.GetIdByOpenId() {
		c.String(404, "服务不存在")
		return
	}
	db := models.Mysql.Exec("update deploy set now_version=now_version+1 where id=?", deployPO.Id)
	if db.Error != nil || db.RowsAffected != 1 {
		c.String(404, "服务不存在")
		return
	}
	relation := &models.DeployServerRelation{
		DeployId: deployPO.Id,
	}
	for _, v := range relation.ListByDeployId() {
		udp_service.Hook001.AddDeploy(v.ServerId)
	}
	c.String(200, "Successful")
}
