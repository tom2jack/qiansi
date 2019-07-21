package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/captcha"
	"qiansi/common/models"
	"qiansi/common/utils"
	"qiansi/common/zmlog"
	"time"
)

var ZM_LOCK = utils.NewLockTable()

// @Summary 获取图片验证码
// @Produce  json
// @Accept  json
// @Success 200 {object} models.ApiResult "{"code":1,"msg":"","data":{"idkey":"ckFbFAcMo7sy7qGyonAd","img":"data:image/png;base64,iVBORw0..."}}"
// @Router /admin/VerifyByImg [get]
func VerifyByImg(c *gin.Context) {
	if ZM_LOCK.IsLock("VerifyImg-ip:"+c.ClientIP(), 3*time.Second) {
		models.NewApiResult(-5, "当前IP数据请求过频，请稍后再试").Json(c)
		return
	}
	idkey, img := captcha.VerifyByImg("")
	models.NewApiResult(1, "请求成功", map[string]string{
		"idkey": idkey,
		"img":   img,
	}).Json(c)
}

// @Summary 获取短信验证码
// @Produce json
// @Accept  multipart/form-data
// @Param phone formData string true "手机号"
// @Param img_idkey formData string true "图片验证码句柄"
// @Param img_code formData string true "图片验证码"
// @Success 200 {object} models.ApiResult "{"code":1,"msg":"发送成功","data":null}"
// @Router /admin/VerifyBySMS [post]
func VerifyBySMS(c *gin.Context) {
	phone := c.PostForm("phone")
	img_idkey := c.PostForm("img_idkey")
	img_code := c.PostForm("img_code")
	if len(phone) != 11 {
		models.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if !captcha.VerifyCheck(img_idkey, img_code) {
		models.NewApiResult(-5, "验证码无效").Json(c)
		return
	}
	if ZM_LOCK.IsLock("phone-ip:"+c.ClientIP(), 15*time.Minute) {
		models.NewApiResult(-5, "当前IP数据请求过频，请稍后再试").Json(c)
		return
	}
	err := captcha.VerifyBySMS(phone)
	if err != nil {
		zmlog.Warn("[短信发送失败]：%s-%s", phone, err.Error())
		models.NewApiResult(-5, "短信发送失败").Json(c)
		return
	}
	models.NewApiResult(1, "发送成功").Json(c)
}
