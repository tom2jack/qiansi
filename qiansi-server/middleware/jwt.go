package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qiansi/common/utils"
	"strconv"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg string = "登陆信息获取失败"
		token := c.GetHeader("LOGIN-KEY")
		if token != "" {
			str, err := utils.ParseToken(token)
			if err != nil {
				print(err.Error())
				msg = "登录验证失败"
			}
			if str != "" {
				c.Set("LOGIN-TOKEN", str)
				uid, _ := strconv.Atoi(str)
				c.Set("UID", uid)
				if uid > 0 {
					msg = ""
				}
			}
		}

		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"msg":  msg,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
