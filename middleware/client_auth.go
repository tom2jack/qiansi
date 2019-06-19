package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// ClientAuth is ClientAuth middleware
func ClientAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		str := c.GetHeader("SERVER-ID")
		server_id, _ := strconv.Atoi(str)
		if server_id > 0 {
			c.Set("SERVER-ID", server_id)
			c.Next()
		} else {
			c.AbortWithStatus(403)
			c.Abort()
			return
		}
	}
}
