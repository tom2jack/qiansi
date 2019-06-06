package admin

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"tools-server/common/utils"
)

var IP_LOCK = utils.NewLockTable()

// @Summary 获取图片验证码
// @Produce  json
// @Accept  json
// @Success 200 {string} json "{"code":1,
// 	"msg":"",
// 	"data":{
// 		"idkey":"ckFbFAcMo7sy7qGyonAd",
// 		"img":"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAAAyCAMAAACJUtIoAAAAP1BMVEUAAAAWWAtfoVSh45aZ245Ymk2U1olRk0aP0YRXmUxtr2Kq7J9HiTyIyn1wsmVpq15golVvsWRusGNMjkGU1olBbKVYAAAAAXRSTlMAQObYZgAAAixJREFUeJzcmOHuqyAMxa0kS/4xmWa+/7veGAWLnEKhsCW3X+YUDz/bUovTf2Du1wDQnPsiF+mHfpOKKriqbXlMpZ+N9jFEhy3LEiFV+GDfR3L5gxOnJjTjqAKEPxiaMEpLo2ak6vBMJ1FEZYWqycyMyBQlk1UzecomDX+wmHS4Ik1Wj99PdVQIdWF4ZSX9z2oAC2ILSn2B6iVz3Ut5XU8ug9vIO87sLVZhVn/UCOa9RB3ecHeyRuo6hOeZQ+VyV6LfSMUPSmQoe263H2hsSEv9wVilKFAuTjDhK8EkLMU9OSx0egDW5YB9Yo4UsMSQVYUSYD1vD2lytUTnZYdnySaSvicEHUikHP1hLRHs88mV0ptfnYvjksUOdWICdM65YqjC5XnOc1Hq+fuvNAlYa0TkFFEKd86QlsuJC1pebkimzCSLIoZkHF9uaMJUuEN/rJAIOQd3MNffVyuT5FpRBxd7gcp3Jh03EiXeqDkUVqqJqs7DUbLj2i9Wkz5gwsuLNXIQKwrsmJ2g+AKUsYi1PuZur270vfZAAWNLc4Sr3n4SCDZFV9moFNlsG6d6vx8MyNIdQ3eoadsiLjZVjgvXvmYo8ILeGmQgVvvWJt84pPM0TlNtNVSlvvI5tIWn3XTzyYN+98U859TQSffHMwXIUwkf9P8M0uUyVjaB6s/EpW+VuX3Kujaq8D5OT2WoPgquDsa7Zc0G4jtUwXp88P2ddfjEOcD8J85G+xcAAP//veIFumaweqUAAAAASUVORK5CYII="
// 	}
// }"
// @Router /admin/verify/VerifyByImg [get]
func VerifyByImg(c *gin.Context) {
	idkey, img := utils.VerifyByImg("")
	utils.Show(c, 1, "", map[string]string{
		"idkey": idkey,
		"img":   img,
	})
}

// @Summary 获取短信验证码
// @Produce  json
// @Accept  json
// @Param phone path string true "手机号"
// @Param img_idkey path string true "图片验证码句柄"
// @Param img_code path string true "图片验证码"
// @Success 200 {string} json "{"code":1,
// 	"msg":"发送成功",
// 	"data":null
// }"
// @Router /admin/verify/VerifyByImg [get]
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
