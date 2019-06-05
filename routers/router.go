package routers

import (
	"github.com/gin-gonic/gin"
	"tools-server/api/admin"
	"tools-server/api/client"
	"tools-server/middleware"
)

var (
	Router *gin.Engine = gin.Default()
)

func LoadRouter() {
	Router.StaticFile("/", "assets/html/index.html")
	Router.Static("/static", "assets/html")

	Router.GET("/api/index", client.ApiIndex)
	Router.GET("/api/ApiRegServer", client.ApiRegServer)

	admin_route := Router.Group("/admin")
	admin_route.Use(middleware.JWT())
	{
		admin_route.GET("/index", admin.AdminIndex)
	}
}
