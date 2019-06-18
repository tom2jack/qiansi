package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tools-server/common/utils"
)

// ClientAuth is ClientAuth middleware
func ClientAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg string = "登陆信息获取失败"
		token := c.GetHeader("LOGIN-KEY")
		if token != "" {
			str := utils.DecrptogAES(token, "2134")
			if str != "" {
				c.Set("LOGIN-TOKEN", str)
				server_id, _ := strconv.Atoi(str)
				c.Set("SERVER_ID", server_id)
				if server_id > 0 {
					msg = ""
				}
			} else {
				msg = "登录验证失败"
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
