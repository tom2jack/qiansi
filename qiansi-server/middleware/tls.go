package middleware

import (
	"gitee.com/zhimiao/qiansi/common/conf"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func TLS() gin.HandlerFunc {
	return func(c *gin.Context) {
		https_listen := conf.S.MustValue("server", "https_listen")
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     https_listen,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
