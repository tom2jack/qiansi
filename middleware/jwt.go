package middleware

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"tools-server/service/utils"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg string
		token := c.GetHeader("LOGIN-KEY")
		if token == "" {
			msg = "登录信息获取失败"
		} else {
			_, err := utils.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					msg = "登录过期"
				default:
					msg = "登录验证失败"
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
