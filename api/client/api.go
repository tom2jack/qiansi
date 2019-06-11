package client

import (
	"github.com/gin-gonic/gin"
	"github.com/lifei6671/gorand"
	"strconv"
	"tools-server/common/utils"
	"tools-server/models"
)

// @Summary 服务器注册
// @Produce  json
// @Accept  json
// @Param uid query string true "用户UID"
// @Param device query string true "客户端设备号"
// @Success 200 {string} json "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /clinet/ApiRegServer [post]
func ApiRegServer(c *gin.Context) {
	uid, _ := strconv.Atoi(c.Query("uid"))
	if !(uid > 0) {
		utils.Show(c, -4, "用户UID非法", nil)
		return
	}
	device := c.Query("device")
	if len(device) != 36 {
		utils.Show(c, -4, "客户端唯一标识号非法", nil)
		return
	}
	var row int
	models.ZM_Mysql.Table("member").Where("uid = ?", uid).Count(&row)
	if row == 0 {
		utils.Show(c, -5, "用户不存在", nil)
		return
	}
	models.ZM_Mysql.Table("server").Where("device_id=?", device).Count(&row)
	if row > 0 {
		utils.Show(c, -5, "设备已存在，请勿重复注册", nil)
		return
	}
	api_secret := string(gorand.KRand(150, gorand.KC_RAND_KIND_ALL))
	server := &models.Server{
		Uid:       uid,
		ApiSecret: api_secret,
		DeviceId:  device,
		Domain:    c.ClientIP(),
	}
	models.ZM_Mysql.Create(server)
	utils.Show(c, 1, "成功", server)
}
