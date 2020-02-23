/**
 * 校验模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019-8-17 18:31:41
 */

package admin

import (
	"gitee.com/zhimiao/qiansi/common/captcha"
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/req"
	"gitee.com/zhimiao/qiansi/resp"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/gin-gonic/gin"
)

var ZM_LOCK = utils.NewLockTable()

// @Summary 获取图片验证码
// @Produce  json
// @Accept  json
// @Success 200 {object} resp.ApiResult "{"code":1,"msg":"","data":{"idkey":"ckFbFAcMo7sy7qGyonAd","img":"data:image/png;base64,iVBORw0..."}}"
// @Router /admin/VerifyByImg [get]
func VerifyByImg(c *gin.Context) {
	if ZM_LOCK.IsLock("VerifyImg-ip:"+c.ClientIP(), 3*time.Second) {
		resp.NewApiResult(-5, "当前IP数据请求过频，请稍后再试").Json(c)
		return
	}
	idkey, img := captcha.VerifyByImg("")
	resp.NewApiResult(1, "请求成功", map[string]string{
		"idkey": idkey,
		"img":   img,
	}).Json(c)
}

// @Summary 获取短信验证码
// @Produce json
// @Accept  application/json
// @Param body body req.VerifyBySMSParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code":1,"msg":"发送成功","data":null}"
// @Router /admin/VerifyBySMS [post]
func VerifyBySMS(c *gin.Context) {
	param := &req.VerifyBySMSParam{}
	if err := c.Bind(param); err != nil {
		resp.NewApiResult(-4, "入参解析失败"+err.Error()).Json(c)
		return
	}
	if len(param.Phone) != 11 {
		resp.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if !captcha.VerifyCheck(param.ImgIdKey, param.ImgCode) {
		resp.NewApiResult(-5, "验证码无效").Json(c)
		return
	}
	if ZM_LOCK.IsLock("phone-ip:"+c.ClientIP(), 15*time.Minute) {
		resp.NewApiResult(-5, "当前IP数据请求过频，请稍后再试").Json(c)
		return
	}
	err := captcha.VerifyBySMS(param.Phone)
	if err != nil {
		logrus.Warnf("[短信发送失败]：%s-%s", param.Phone, err.Error())
		resp.NewApiResult(-5, "短信发送失败").Json(c)
		return
	}
	resp.NewApiResult(1, "发送成功").Json(c)
}
