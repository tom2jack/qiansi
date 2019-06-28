package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/utils"
	"qiansi/common/zmlog"
	"time"
)

var ZM_LOCK = utils.NewLockTable()

// @Summary 获取图片验证码
// @Produce  json
// @Accept  json
// @Success 200 {object} utils.Json "{"code":1,"msg":"","data":{"idkey":"ckFbFAcMo7sy7qGyonAd","img":"data:image/png;base64,iVBORw0..."}}"
// @Router /admin/VerifyByImg [get]
func VerifyByImg(c *gin.Context) {
	if ZM_LOCK.IsLock("VerifyImg-ip:"+c.ClientIP(), 3*time.Second) {
		utils.Show(c, -5, "当前IP数据请求过频，请稍后再试", nil)
		return
	}
	idkey, img := utils.VerifyByImg("")
	utils.Show(c, 1, "请求成功", map[string]string{
		"idkey": idkey,
		"img":   img,
	})
}

// @Summary 获取短信验证码
// @Produce json
// @Accept  multipart/form-data
// @Param phone formData string true "手机号"
// @Param img_idkey formData string true "图片验证码句柄"
// @Param img_code formData string true "图片验证码"
// @Success 200 {object} utils.Json "{"code":1,"msg":"发送成功","data":null}"
// @Router /admin/VerifyBySMS [post]
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
	if ZM_LOCK.IsLock("phone-ip:"+c.ClientIP(), 15*time.Minute) {
		utils.Show(c, -5, "当前IP数据请求过频，请稍后再试", nil)
		return
	}
	err := utils.VerifyBySMS(phone)
	if err != nil {
		zmlog.Warn("[短信发送失败]：%s-%s", phone, err.Error())
		utils.Show(c, -5, "短信发送失败", nil)
		return
	}
	utils.Show(c, 1, "发送成功", nil)
}
