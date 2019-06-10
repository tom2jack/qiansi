package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"tools-server/conf"
)

func TLS() gin.HandlerFunc {
	return func(c *gin.Context) {
		endPoint := conf.App.MustValue("server", "http_listen", ":7091")
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     endPoint,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
