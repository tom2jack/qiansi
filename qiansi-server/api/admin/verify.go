package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/captcha"
	"qiansi/common/models"
	"qiansi/common/models/api_req"
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
// @Accept  application/json
// @Param body body api_req.VerifyBySMSParam true "入参集合"
// @Success 200 {object} models.ApiResult "{"code":1,"msg":"发送成功","data":null}"
// @Router /admin/VerifyBySMS [post]
func VerifyBySMS(c *gin.Context) {
	param := &api_req.VerifyBySMSParam{}
	if err := c.Bind(param); err != nil {
		models.NewApiResult(-4, "入参解析失败"+err.Error()).Json(c)
		return
	}
	if len(param.Phone) != 11 {
		models.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if !captcha.VerifyCheck(param.ImgIdKey, param.ImgCode) {
		models.NewApiResult(-5, "验证码无效").Json(c)
		return
	}
	if ZM_LOCK.IsLock("phone-ip:"+c.ClientIP(), 15*time.Minute) {
		models.NewApiResult(-5, "当前IP数据请求过频，请稍后再试").Json(c)
		return
	}
	err := captcha.VerifyBySMS(param.Phone)
	if err != nil {
		zmlog.Warn("[短信发送失败]：%s-%s", param.Phone, err.Error())
		models.NewApiResult(-5, "短信发送失败").Json(c)
		return
	}
	models.NewApiResult(1, "发送成功").Json(c)
}
