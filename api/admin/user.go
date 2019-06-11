package admin

import (
	"github.com/gin-gonic/gin"
	"time"
	"tools-server/common/utils"
	"tools-server/models"
)

// @Summary 登录
// @Produce  json
// @Accept  multipart/form-data
// @Param phone query string true "手机号"
// @Param password query string true "密码"
// @Success 200 {string} json "{"code": 1,"msg": "登录成功","data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00"}}"
// @Router /admin/user/UserSigin [post]
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
	if IP_LOCK.IsLock("user-sigin:"+phone, 5*time.Second) {
		utils.Show(c, -5, "操作过频", nil)
		return
	}
	member := &models.Member{}
	models.ZM_Mysql.Table("member").Where("phone = ?", phone).First(member)
	if !utils.PasswordVerify(member.Password, password) {
		utils.Show(c, -5, "密码错误", nil)
		return
	}
	member.Password = ""
	utils.Show(c, 1, "登录成功", member)
}

// @Summary 注册账号
// @Produce  json
// @Accept  multipart/form-data
// @Param phone query string true "手机号"
// @Param password query string true "密码"
// @Success 200 {string} json "{"code": 1,"msg": "登录成功","data": {"CreateTime": "2019-02-27T16:11:27+08:00","InviterUid": 0,"Password": "","Phone": "15061370322","Status": 1,"Uid": 2, "UpdateTime": "2019-02-27T16:19:54+08:00"}}"
// @Router /admin/user/UserSiginUp [post]
func UserSiginUp(c *gin.Context) {
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
	if IP_LOCK.IsLock("user-sigin:"+phone, 5*time.Second) {
		utils.Show(c, -5, "操作过频", nil)
		return
	}
	member := &models.Member{}
	models.ZM_Mysql.Table("member").Where("phone = ?", phone).First(member)
	if !utils.PasswordVerify(member.Password, password) {
		utils.Show(c, -5, "密码错误", nil)
		return
	}
	member.Password = ""
	utils.Show(c, 1, "登录成功", member)
}
