package client

import (
	"github.com/gin-gonic/gin"
	"github.com/lifei6671/gorand"
	"strconv"
	"tools-server/common/utils"
	"tools-server/models"
)

func ApiIndex(c *gin.Context) {
	utils.VerifyBySMS("15061370322")
	idkey, img := utils.VerifyByImg("")
	utils.Show(c, 1, "", map[string]string{
		"idkey": idkey,
		"img":   img,
	})
}

//ApiRegServer 服务器注册
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
