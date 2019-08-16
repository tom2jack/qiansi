package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// ClientAuth is ClientAuth middleware
func ClientAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverId, _ := strconv.Atoi(c.GetHeader("SERVER-ID"))
		serverUid, _ := strconv.Atoi(c.GetHeader("SERVER-UID"))
		serverDevice := c.GetHeader("SERVER-DEVICE")
		if serverId == 0 || serverUid == 0 || serverDevice == "" {
			c.AbortWithStatus(403)
			c.Abort()
			return
		}
		c.Set("SERVER-ID", serverId)
		c.Set("SERVER-UID", serverUid)
		c.Set("SERVER-DEVICE", serverDevice)
		c.Next()
	}
}
