package api

import (
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/config"
	"github.com/zhi-miao/qiansi/common/errors"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"

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
		msg := "登陆信息获取失败"
		token := c.GetHeader(req.LOGIN_KEY)
		if token != "" {
			str, err := gutils.ParseToken(token, []byte(config.GetConfig().App.JwtSecret))
			if err != nil || str == "" {
				msg = "登录验证失败"
			} else {
				uid, _ := strconv.Atoi(str)
				if uid > 0 {
					msg = ""
					c.Set(req.LOGIN_TOKEN, str)
					c.Set(req.UID, uid)
				}
			}
		}
		if msg != "" {
			c.JSON(resp.ApiError(errors.Unauthorized, msg))
			c.Abort()
			return
		}
		c.Next()
	}
}
