package admin

import (
	"github.com/gin-gonic/gin"
	"qiansi/common/captcha"
	"qiansi/common/models"
	"qiansi/common/utils"
	"strconv"
	"time"
)

type UserInfo struct {
	models.Member
	Token string
}

// @Summary 登录
// @Produce  json
// @Accept  multipart/form-data
// @Param phone formData string true "手机号"
// @Param password formData string true "密码"
// @Success 200 {object} models.ApiResult "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSigin [post]
func UserSigin(c *gin.Context) {
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	if len(phone) != 11 {
		models.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if len(password) < 6 {
		models.NewApiResult(-4, "密码错误").Json(c)
		return
	}
	if ZM_LOCK.IsLock("user-sigin:"+phone, 5*time.Second) {
		models.NewApiResult(-5, "操作过频").Json(c)
		return
	}
	member := &models.Member{}
	models.ZM_Mysql.Table("member").Where("phone = ?", phone).First(member)
	if !utils.PasswordVerify(member.Password, password) {
		models.NewApiResult(-5, "密码错误，不给登陆").Json(c)
		return
	}
	member.Password = ""
	token, err := utils.CreateToken(strconv.Itoa(member.Uid), 24*time.Hour)
	if err != nil {
		models.NewApiResult(-5, "Token生成失败，无法登陆，请联系管理员").Json(c)
		return
	}
	u := &UserInfo{
		*member,
		token,
	}
	models.NewApiResult(1, "登陆成功", u).Json(c)
}

// @Summary 注册账号
// @Produce  json
// @Accept  multipart/form-data
// @Param phone formData string true "手机号"
// @Param password formData string true "密码"
// @Param code formData string true "短信验证码"
// @Param inviter_uid formData int false "邀请人UID"
// @Success 200 {object} models.ApiResult "{"code": 1,"msg": "注册成功","data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSiginUp [post]
func UserSiginUp(c *gin.Context) {
	var row int
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	code := c.PostForm("code")
	inviter_uid, _ := strconv.Atoi(c.PostForm("inviter_uid"))
	if len(phone) != 11 {
		models.NewApiResult(-4, "手机号错误").Json(c)
		return
	}
	if len(password) < 6 {
		models.NewApiResult(-4, "密码格式不正确").Json(c)
		return
	}
	if ZM_LOCK.IsLock("user-siginup:"+phone, 15*time.Second) {
		models.NewApiResult(-5, "操作过频").Json(c)
		return
	}
	// 验证邀请人
	if inviter_uid > 0 {
		models.ZM_Mysql.Table("member").Where("uid=? and status=1", inviter_uid).Count(&row)
		if row > 0 {
			models.NewApiResult(-5, "非法的邀请人").Json(c)
			return
		}
	}
	if len(code) < 4 || !captcha.VerifyCheck(captcha.VerifyBySMSIDKEY(phone), code) {
		models.NewApiResult(-4, "验证码错误").Json(c)
		return
	}
	models.ZM_Mysql.Table("member").Where("phone=?", phone).Count(&row)
	if row > 0 {
		models.NewApiResult(-5, "请勿重复注册").Json(c)
		return
	}
	member := &models.Member{
		Phone:      phone,
		Password:   utils.PasswordHash(password),
		InviterUid: inviter_uid,
		Status:     1,
	}
	models.ZM_Mysql.Create(member)
	member.Password = ""
	token, err := utils.CreateToken(strconv.Itoa(member.Uid), 24*time.Hour)
	if err != nil {
		models.NewApiResult(-5, "Token生成失败，无法登陆，请联系管理员").Json(c)
		return
	}
	u := &UserInfo{
		*member,
		token,
	}
	if member.Uid > 0 {
		models.NewApiResult(1, "注册成功", u).Json(c)
		return
	}
	models.NewApiResult(0, "注册失败", u).Json(c)
}

// @Summary 修改密码
// @Produce  json
// @Accept  multipart/form-data
// @Param old_password formData string true "旧密码"
// @Param new_password formData string true "新密码"
// @Success 200 {object} models.ApiResult "{"code": 1,"msg": "修改成功", "data": null}}"
// @Router /admin/UserResetPwd [post]
func UserResetPwd(c *gin.Context) {
	old_password := c.PostForm("old_password")
	new_password := c.PostForm("new_password")
	if len(old_password) < 6 || len(new_password) < 6 {
		models.NewApiResult(-4, "密码格式错误").Json(c)
		return
	}
	if old_password == new_password {
		models.NewApiResult(-4, "新旧密码不能一样").Json(c)
		return
	}
	if ZM_LOCK.IsLock("UserResetPwd:"+c.GetString("LOGIN-TOKEN"), 30*time.Hour) {
		models.NewApiResult(-5, "密码每30分钟只能尝试修改一次").Json(c)
		return
	}
	member := &models.Member{}
	models.ZM_Mysql.Table("member").Where("uid = ?", c.GetInt("UID")).First(member)
	if !utils.PasswordVerify(member.Password, old_password) {
		models.NewApiResult(-5, "密码错误").Json(c)
		return
	}
	models.ZM_Mysql.Table("member").Where("uid = ?", c.GetInt("UID")).Update("password", utils.PasswordHash(new_password))
	models.NewApiResult(1, "修改成功").Json(c)
}
