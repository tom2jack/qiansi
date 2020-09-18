package admin

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/captcha"
	"github.com/zhi-miao/qiansi/common/errors"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"

	"github.com/gin-gonic/gin"
)

type verifyApi struct{}

var Verify = &verifyApi{}

var ZM_LOCK = gutils.NewLockTable()

// @Summary 获取图片验证码
// @Produce  json
// @Accept  json
// @Success 200 {object} gin.H "{"idkey":"ckFbFAcMo7sy7qGyonAd","img":"data:image/png;base64,iVBORw0..."}"
// @Router /admin/VerifyByImg [get]
func (r *verifyApi) ByImg(c *gin.Context) {
	if ZM_LOCK.IsLock("VerifyImg-ip:"+c.ClientIP(), 3*time.Second) {
		c.JSON(resp.ApiError("当前IP数据请求过频，请稍后再试"))
		return
	}
	idkey, img := captcha.VerifyByImg("")
	c.JSON(resp.ApiSuccess(gin.H{
		"idkey": idkey,
		"img":   img,
	}))
}

// @Summary 获取短信验证码
// @Produce json
// @Accept  application/json
// @Param body body req.VerifyBySMSParam true "入参集合"
// @Success 200
// @Router /admin/VerifyBySMS [post]
func (r *verifyApi) BySMS(c *gin.Context) {
	param := &req.VerifyBySMSParam{}
	if err := c.Bind(param); err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
	if len(param.Phone) != 11 {
		c.JSON(resp.ApiError("手机号错误"))
		return
	}
	if !captcha.VerifyCheck(param.ImgIdKey, param.ImgCode) {
		c.JSON(resp.ApiError("验证码无效"))
		return
	}
	if ZM_LOCK.IsLock("phone-ip:"+c.ClientIP(), 15*time.Minute) {
		c.JSON(resp.ApiError("当前IP数据请求过频，请稍后再试"))
		return
	}
	err := captcha.VerifyBySMS(param.Phone)
	if err != nil {
		logrus.Warnf("[短信发送失败]：%s-%s", param.Phone, err.Error())
		c.JSON(resp.ApiError(errors.InternalServerError, "短信发送失败"))
		return
	}
}
