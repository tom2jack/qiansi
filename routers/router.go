package routers

import (
	"github.com/gin-gonic/gin"
	"tools-server/api/admin"
	"tools-server/api/api"
	"tools-server/middleware"
)

var (
	Router *gin.Engine = gin.Default()
)

func init() {
	Router.GET("/api/index", api.ApiIndex)
	Router.StaticFile("/", "assets/html/index.html")
	Router.Static("/static", "assets/html")

	admin_route := Router.Group("/admin")
	admin_route.Use(middleware.JWT())
	{
		admin_route.GET("/index", admin.AdminIndex)
	}
}
