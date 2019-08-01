package client

import (
	"github.com/gin-gonic/gin"
	"github.com/lifei6671/gorand"
	"qiansi/common/models"
	"qiansi/common/utils"
	"qiansi/common/zmlog"
	"qiansi/qiansi-server/net_service/udp_service"
	"strconv"
)

// @Summary 服务器注册
// @Produce  json
// @Accept  json
// @Param uid query string true "用户UID"
// @Param device query string true "客户端设备号"
// @Success 200 {object} models.ApiResult "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /clinet/ApiRegServer [post]
func ApiRegServer(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("uid"))
	if !(uid > 0) {
		models.NewApiResult(-4, "用户UID非法").Json(c)
		return
	}
	device := c.Query("device")
	if len(device) != 36 {
		models.NewApiResult(-4, "客户端唯一标识号非法").Json(c)
		return
	}
	var row int
	models.ZM_Mysql.Table("member").Where("uid = ?", uid).Count(&row)
	if row == 0 {
		models.NewApiResult(-5, "用户不存在").Json(c)
		return
	}
	models.ZM_Mysql.Table("server").Where("device_id=?", device).Count(&row)
	if row > 0 {
		models.NewApiResult(-5, "设备已存在，请勿重复注册").Json(c)
		return
	}
	api_secret := string(gorand.KRand(16, gorand.KC_RAND_KIND_ALL))
	server := &models.Server{
		Uid:       uid,
		ApiSecret: api_secret,
		DeviceId:  device,
		Domain:    c.ClientIP(),
	}
	models.ZM_Mysql.Create(server)
	models.NewApiResult(1, "成功", server).Json(c)
}

// @Summary 获取服务器部署任务清单
// @Produce  json
// @Accept  json
// @Success 200 {object} models.ApiResult "{"code": 1,"msg": "读取成功","data": [deploy]}"
// @Router /clinet/ApiGetDeployTask [GET]
func ApiGetDeployTask(c *gin.Context) {
	server_id := c.GetInt("SERVER-ID")
	defer udp_service.Hook001.Deploy.DEL(strconv.Itoa(server_id))
	deploy := &[]models.Deploy{}
	// models.ZM_Mysql.Raw()
	models.ZM_Mysql.Raw("SELECT d.* FROM `deploy` d LEFT JOIN `deploy_server_relation` r ON d.id=r.deploy_id WHERE r.server_id=? and d.now_version > r.deploy_version", server_id).Scan(deploy)
	models.NewApiResult(1, "读取成功", deploy).Encypt(c)
}

// @Summary 客户端日志推送
// @Produce  json
// @Accept  multipart/form-data
// @Param server_id formData string true "客户端平台编号"
// @Param device formData string true "客户端设备号"
// @Param content formData string true "日志文本内容"
// @Success 200 {object} models.ApiResult ""
// @Router /clinet/LogPush [post]
func LogPush(c *gin.Context) {
	serverId, _ := strconv.Atoi(c.PostForm("server_id"))
	device := c.PostForm("device")
	content := utils.MustUtf8(c.PostForm("content"))
	var row int
	models.ZM_Mysql.Table("server").Where("id=? and device_id=?", serverId, device).Count(&row)
	if row == 0 {
		c.Status(403)
		return
	}
	serverLog := &models.ServerLog{
		ServerId: serverId,
		DeviceId: device,
		Content:  content,
	}
	models.ZM_Mysql.Create(serverLog)
	models.NewApiResult(1).Json(c)
}

// @Summary 客户端部署成功回调
// @Produce  json
// @Accept  multipart/form-data
// @Param server_id formData string true "客户端平台编号"
// @Param device formData string true "客户端设备号"
// @Param deploy_id formData string true "日志文本内容"
// @Success 200 {object} models.ApiResult ""
// @Router /clinet/DeployNotify [post]
func DeployNotify(c *gin.Context) {
	// TODO: 实现版本号的增长，用于剔除下次部署发布中此机器的版本号
	serverId, _ := strconv.Atoi(c.PostForm("server_id"))
	device := c.PostForm("device")
	zmlog.Info("%d,%s", serverId, device)
}
