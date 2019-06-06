package admin

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"tools-server/common/utils"
)

var IP_LOCK = utils.NewLockTable()

func VerifyByImg(c *gin.Context) {
	idkey, img := utils.VerifyByImg("")
	utils.Show(c, 1, "", map[string]string{
		"idkey": idkey,
		"img":   img,
	})
}

func VerifyBySMS(c *gin.Context) {
	phone := c.PostForm("phone")
	img_idkey := c.PostForm("img_idkey")
	img_code := c.PostForm("img_code")
	if len(phone) != 11 {
		utils.Show(c, -4, "手机号错误", nil)
		return
	}
	if !utils.VerifyCheck(img_idkey, img_code) {
		utils.Show(c, -5, "验证码无效", nil)
		return
	}
	if IP_LOCK.IsLock("phone-ip:"+c.ClientIP(), 15*time.Minute) {
		utils.Show(c, -5, "当前IP数据请求过频，请稍后再试", nil)
		return
	}
	err := utils.VerifyBySMS(phone)
	if err != nil {
		log.Printf("[短信发送失败]：%s-%s", phone, err.Error())
		utils.Show(c, -5, "短信发送失败", nil)
		return
	}
	utils.Show(c, 1, "发送成功", nil)
}
