package api

import (
	"gitee.com/zhimiao/qiansi/common/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// corsMiddleware 跨域
func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	})
}

// logMiddleware 日志中间件
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUrl := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求ip
		clientIP := c.ClientIP()
		// 日志格式
		logrus.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}

// clientAuthMiddleware 客户端认证中间件
func clientAuthMiddleware() gin.HandlerFunc {
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

// jwtMiddleware jwt鉴权
func jwtMiddleware() gin.HandlerFunc {
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
