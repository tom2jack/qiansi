/**
 * 用户模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019-8-17 18:31:41
 */

package admin

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/captcha"
	"github.com/zhi-miao/qiansi/common/config"
	"github.com/zhi-miao/qiansi/common/errors"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/service"

	"github.com/gin-gonic/gin"
)

type userApi struct{}

var User = &userApi{}

// @Summary 登录
// @Produce  json
// @Accept json
// @Param body body req.UserSiginParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSigin [post]
func (r *userApi) Sigin(c *gin.Context) {
	param := &req.UserSiginParam{}
	if err := c.Bind(param); err != nil {
		c.JSON(resp.ApiError("入参解析失败"))
		return
	}
	if len(param.Phone) != 11 {
		c.JSON(resp.ApiError("手机号错误"))
		return
	}
	if len(param.Password) < 6 {
		c.JSON(resp.ApiError("密码错误"))
		return
	}
	if ZM_LOCK.IsLock("user-sigin:"+param.Phone, 5*time.Second) {
		c.JSON(resp.ApiError("操作过频"))
		return
	}
	member := &models.Member{}
	models.Mysql.Table("member").Where("phone = ?", param.Phone).First(member)
	if !gutils.PasswordVerify(member.Password, param.Password) {
		c.JSON(resp.ApiError("密码错误，不给登陆"))
		return
	}
	member.Password = ""
	token, err := gutils.CreateToken(strconv.Itoa(member.Id), 9*24*time.Hour, []byte(config.GetConfig().App.JwtSecret))
	if err != nil {
		c.JSON(resp.ApiError("Token生成失败，无法登陆，请联系管理员"))
		return
	}
	c.JSON(resp.ApiSuccess(&resp.UserInfoVO{
		Member: *member,
		Token:  token,
	}))
}

// @Summary 注册账号
// @Produce  json
// @Accept json
// @Param body body req.UserSiginUpParam true "入参集合"
// @Success 200 {object} resp.UserInfoVO "响应成功"
// @Success 400 {object} resp.ApiErrorResult "错误"
// @Router /admin/UserSiginUp [post]
func (r *userApi) SiginUp(c *gin.Context) {
	param := &req.UserSiginUpParam{}
	if err := c.Bind(param); err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
	if len(param.Phone) != 11 {
		c.JSON(resp.ApiError("手机号错误"))
		return
	}
	if len(param.Password) < 6 {
		c.JSON(resp.ApiError("密码格式不正确"))
		return
	}
	if ZM_LOCK.IsLock("user-siginup:"+param.Phone, 15*time.Second) {
		c.JSON(resp.ApiError("操作过频"))
		return
	}
	userModel := models.GetMemberModels()
	// 验证邀请人
	if param.InviterUid > 0 {
		if !userModel.ExistsUID(param.InviterUid) {
			c.JSON(resp.ApiError("非法的邀请人"))
			return
		}
	}
	if len(param.Code) < 4 || !captcha.VerifyCheck("phone:"+param.Phone, param.Code) {
		c.JSON(resp.ApiError("验证码错误"))
		return
	}
	if !userModel.ExistsPhone(param.Phone) {
		c.JSON(resp.ApiError("请勿重复注册"))
		return
	}
	member := &models.Member{
		Phone:       param.Phone,
		Password:    gutils.PasswordHash(param.Password),
		InviterUid:  param.InviterUid,
		MaxSchedule: 2,
		MaxDeploy:   2,
	}
	models.Mysql.Table("member").Create(member)
	member.Password = ""
	token, err := gutils.CreateToken(strconv.Itoa(member.Id), 24*time.Hour, []byte(config.GetConfig().App.JwtSecret))
	if err != nil {
		c.JSON(resp.ApiError("Token生成失败，无法登陆，请联系管理员"))
		return
	}
	if member.Id == 0 {
		c.JSON(resp.ApiError(errors.InternalServerError, "注册失败"))
		return
	}
	// 用户邀请奖励发放
	err = service.GetActivityService().Inviter(param.InviterUid, member.Id)
	if err != nil {
		logrus.Warnf("邀请奖励%d->%d发放失败, %s", member.Id, param.InviterUid, err.Error())
	}
	c.JSON(resp.ApiSuccess(&resp.UserInfoVO{
		Member: *member,
		Token:  token,
	}))
}

// @Summary 修改密码
// @Produce  json
// @Accept  json
// @Param body body req.UserResetPwdParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "修改成功", "data": null}"
// @Router /admin/UserResetPwd [post]
func (r *userApi) ResetPwd(c *gin.Context) {
	param := &req.UserResetPwdParam{}
	if err := c.Bind(param); err != nil {
		c.JSON(resp.ApiError(err))
		return
	}
	if len(param.OldPassword) < 6 || len(param.NewPassword) < 6 {
		c.JSON(resp.ApiError("密码格式错误"))
		return
	}
	if param.OldPassword == param.NewPassword {
		c.JSON(resp.ApiError("新旧密码不能一样"))
		return
	}
	if ZM_LOCK.IsLock("UserResetPwd:"+c.GetString("LOGIN-TOKEN"), 30*time.Hour) {
		c.JSON(resp.ApiError("密码每30分钟只能尝试修改一次"))
		return
	}
	member := &models.Member{}
	models.Mysql.Table("member").Where("id = ?", c.GetInt("UID")).First(member)
	if !gutils.PasswordVerify(member.Password, param.OldPassword) {
		c.JSON(resp.ApiError("密码错误"))
		return
	}
	models.Mysql.Table("member").Where("id = ?", c.GetInt("UID")).Update("password", gutils.PasswordHash(param.NewPassword))
}
