/**
 * 用户模块
 * Created by 纸喵软件.
 * User: 倒霉狐狸
 * Date: 2019-8-17 18:31:41
 */

package admin

import (
	"qiansi/common/captcha"
	"qiansi/common/utils"
	"qiansi/qiansi-server/models"
	"qiansi/qiansi-server/req"
	"qiansi/qiansi-server/resp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 登录
// @Produce  json
// @Accept json
// @Param body body req.UserSiginParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSigin [post]
func UserSigin(c *gin.Context) {
	param := &req.UserSiginParam{}
	if err := c.Bind(param); err != nil {
		resp.NewApiResult(-4, "入参解析失败").Json(c)
		return
	}
	if len(param.Phone) != 11 {
		resp.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if len(param.Password) < 6 {
		resp.NewApiResult(-4, "密码错误").Json(c)
		return
	}
	if ZM_LOCK.IsLock("user-sigin:"+param.Phone, 5*time.Second) {
		resp.NewApiResult(-5, "操作过频").Json(c)
		return
	}
	member := &models.Member{}
	models.Mysql.Table("member").Where("phone = ?", param.Phone).First(member)
	if !utils.PasswordVerify(member.Password, param.Password) {
		resp.NewApiResult(-5, "密码错误，不给登陆").Json(c)
		return
	}
	member.Password = ""
	token, err := utils.CreateToken(strconv.Itoa(member.Id), 9*24*time.Hour)
	if err != nil {
		resp.NewApiResult(-5, "Token生成失败，无法登陆，请联系管理员").Json(c)
		return
	}
	resp.NewApiResult(1, "登陆成功", &resp.UserInfoVO{
		*member,
		token,
	}).Json(c)
}

// @Summary 注册账号
// @Produce  json
// @Accept json
// @Param body body req.UserSiginUpParam true "入参集合"
// @Success 200 {object} resp.UserInfoVO "{"code": 1,"msg": "注册成功","data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSiginUp [post]
func UserSiginUp(c *gin.Context) {
	var row int
	param := &req.UserSiginUpParam{}
	if err := c.Bind(param); err != nil {
		resp.NewApiResult(-4, "入参解析失败").Json(c)
		return
	}
	if len(param.Phone) != 11 {
		resp.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if len(param.Password) < 6 {
		resp.NewApiResult(-4, "密码格式不正确").Json(c)
		return
	}
	if ZM_LOCK.IsLock("user-siginup:"+param.Phone, 15*time.Second) {
		resp.NewApiResult(-5, "操作过频").Json(c)
		return
	}
	// 验证邀请人
	if param.InviterUid > 0 {
		models.Mysql.Table("member").Where("uid=? and status=1", param.InviterUid).Count(&row)
		if row > 0 {
			resp.NewApiResult(-5, "非法的邀请人").Json(c)
			return
		}
	}
	if len(param.Code) < 4 || !captcha.VerifyCheck("phone:"+param.Phone, param.Code) {
		resp.NewApiResult(-4, "验证码错误").Json(c)
		return
	}
	models.Mysql.Table("member").Where("phone=?", param.Phone).Count(&row)
	if row > 0 {
		resp.NewApiResult(-5, "请勿重复注册").Json(c)
		return
	}
	member := &models.Member{
		Phone:      param.Phone,
		Password:   utils.PasswordHash(param.Password),
		InviterUid: param.InviterUid,
	}
	models.Mysql.Table("member").Create(member)
	member.Password = ""
	token, err := utils.CreateToken(strconv.Itoa(member.Id), 24*time.Hour)
	if err != nil {
		resp.NewApiResult(-5, "Token生成失败，无法登陆，请联系管理员").Json(c)
		return
	}
	if member.Id > 0 {
		resp.NewApiResult(1, "注册成功", &resp.UserInfoVO{
			*member,
			token,
		}).Json(c)
		return
	}
	resp.NewApiResult(0, "注册失败").Json(c)
}

// @Summary 修改密码
// @Produce  json
// @Accept  json
// @Param body body req.UserResetPwdParam true "入参集合"
// @Success 200 {object} resp.ApiResult "{"code": 1,"msg": "修改成功", "data": null}"
// @Router /admin/UserResetPwd [post]
func UserResetPwd(c *gin.Context) {
	param := &req.UserResetPwdParam{}
	if err := c.Bind(param); err != nil {
		resp.NewApiResult(-4, "入参解析失败").Json(c)
		return
	}
	if len(param.OldPassword) < 6 || len(param.NewPassword) < 6 {
		resp.NewApiResult(-4, "密码格式错误").Json(c)
		return
	}
	if param.OldPassword == param.NewPassword {
		resp.NewApiResult(-4, "新旧密码不能一样").Json(c)
		return
	}
	if ZM_LOCK.IsLock("UserResetPwd:"+c.GetString("LOGIN-TOKEN"), 30*time.Hour) {
		resp.NewApiResult(-5, "密码每30分钟只能尝试修改一次").Json(c)
		return
	}
	member := &models.Member{}
	models.Mysql.Table("member").Where("id = ?", c.GetInt("UID")).First(member)
	if !utils.PasswordVerify(member.Password, param.OldPassword) {
		resp.NewApiResult(-5, "密码错误").Json(c)
		return
	}
	models.Mysql.Table("member").Where("id = ?", c.GetInt("UID")).Update("password", utils.PasswordHash(param.NewPassword))
	resp.NewApiResult(1, "修改成功").Json(c)
}
