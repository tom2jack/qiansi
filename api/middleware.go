package api

import (
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/config"
	"github.com/zhi-miao/qiansi/common/req"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

// jwtMiddleware jwt鉴权
func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg string = "登陆信息获取失败"
		token := c.GetHeader("LOGIN-KEY")
		if token != "" {
			str, err := gutils.ParseToken(token, []byte(config.GetConfig().App.JwtSecret))
			if err != nil {
				print(err.Error())
				msg = "登录验证失败"
			}
			if str != "" {
				c.Set("LOGIN-TOKEN", str)
				uid, _ := strconv.Atoi(str)
				c.Set(req.UID, uid)
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
