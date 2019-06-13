package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"tools-server/common/utils"
	"tools-server/models"
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
// @Success 200 {object} utils.Json "{"code": 1,"msg": "登录成功", "data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSigin [post]
func UserSigin(c *gin.Context) {
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	if len(phone) != 11 {
		utils.Show(c, -4, "手机号错误", nil)
		return
	}
	if len(password) < 6 {
		utils.Show(c, -4, "密码错误", nil)
		return
	}
	if ZM_LOCK.IsLock("user-sigin:"+phone, 5*time.Second) {
		utils.Show(c, -5, "操作过频", nil)
		return
	}
	member := &models.Member{}
	models.ZM_Mysql.Table("member").Where("phone = ?", phone).First(member)
	if !utils.PasswordVerify(member.Password, password) {
		utils.Show(c, -5, "密码错误，不给登陆", nil)
		return
	}
	member.Password = ""
	token, err := utils.CreateToken(strconv.Itoa(member.Uid), 24*time.Hour)
	if err != nil {
		utils.Show(c, -5, "Token生成失败，无法登陆，请联系管理员", nil)
		return
	}
	u := &UserInfo{
		*member,
		token,
	}
	utils.Show(c, 1, "登陆成功", u)
}

// @Summary 注册账号
// @Produce  json
// @Accept  multipart/form-data
// @Param phone formData string true "手机号"
// @Param password formData string true "密码"
// @Param code formData string true "短信验证码"
// @Param inviter_uid formData int false "邀请人UID"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "注册成功","data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00", "Token":"sdfsdafsd.."}}"
// @Router /admin/UserSiginUp [post]
func UserSiginUp(c *gin.Context) {
	var row int
	phone := c.PostForm("phone")
	password := c.PostForm("password")
	code := c.PostForm("code")
	inviter_uid, _ := strconv.Atoi(c.PostForm("inviter_uid"))
	if len(phone) != 11 {
		utils.Show(c, -4, "手机号错误", nil)
		return
	}
	if len(password) < 6 {
		utils.Show(c, -4, "密码格式不正确", nil)
		return
	}
	if ZM_LOCK.IsLock("user-siginup:"+phone, 15*time.Second) {
		utils.Show(c, -5, "操作过频", nil)
		return
	}
	// 验证邀请人
	if inviter_uid > 0 {
		models.ZM_Mysql.Table("member").Where("uid=? and status=1", inviter_uid).Count(&row)
		if row > 0 {
			utils.Show(c, -5, "非法的邀请人", nil)
			return
		}
	}
	if len(code) < 4 || !utils.VerifyCheck(utils.VerifyBySMSIDKEY(phone), code) {
		utils.Show(c, -4, "验证码错误", nil)
		return
	}
	models.ZM_Mysql.Table("member").Where("phone=?", phone).Count(&row)
	if row > 0 {
		utils.Show(c, -5, "请勿重复注册", nil)
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
		utils.Show(c, -5, "Token生成失败，无法登陆，请联系管理员", nil)
		return
	}
	u := &UserInfo{
		*member,
		token,
	}
	if member.Uid > 0 {
		utils.Show(c, 1, "注册成功", u)
		return
	}
	utils.Show(c, 0, "注册失败", u)
}

// @Summary 修改密码
// @Produce  json
// @Accept  multipart/form-data
// @Param old_password formData string true "旧密码"
// @Param new_password formData string true "新密码"
// @Success 200 {object} utils.Json "{"code": 1,"msg": "修改成功", "data": null}}"
// @Router /admin/UserResetPwd [post]
func UserResetPwd(c *gin.Context) {
	old_password := c.PostForm("old_password")
	new_password := c.PostForm("new_password")
	if len(old_password) < 6 || len(new_password) < 6 {
		utils.Show(c, -4, "密码格式错误", nil)
		return
	}
	if old_password == new_password {
		utils.Show(c, -4, "新旧密码不能一样", nil)
		return
	}
	if ZM_LOCK.IsLock("UserResetPwd:"+c.GetString("LOGIN-TOKEN"), 30*time.Hour) {
		utils.Show(c, -5, "密码每30分钟只能尝试修改一次", nil)
		return
	}
	member := &models.Member{}
	models.ZM_Mysql.Table("member").Where("uid = ?", c.GetInt("UID")).First(member)
	if !utils.PasswordVerify(member.Password, old_password) {
		utils.Show(c, -5, "密码错误", nil)
		return
	}
	models.ZM_Mysql.Table("member").Where("uid = ?", c.GetInt("UID")).Update("password", utils.PasswordHash(new_password))
	utils.Show(c, 1, "修改成功", nil)
}
