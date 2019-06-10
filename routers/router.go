package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"tools-server/api/admin"
	"tools-server/api/client"
	_ "tools-server/docs"
	"tools-server/middleware"
)

var (
	Router *gin.Engine = gin.Default()
)

func LoadRouter() {
	// 加载https跳转
	Router.Use(middleware.TLS())

	// use ginSwagger middleware to
	Router.GET("/docs/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	Router.StaticFile("/", "assets/html/index.html")
	Router.Static("/static", "assets/html")

	Router.GET("/api/index", client.ApiIndex)
	Router.GET("/api/ApiRegServer", client.ApiRegServer)

	admin_route := Router.Group("/admin")
	admin_route.GET("/verify/VerifyByImg", admin.VerifyByImg)
	admin_route.POST("/verify/VerifyBySMS", admin.VerifyBySMS)
	admin_route.POST("/user/UserSigin", admin.UserSigin)
	admin_route.Use(middleware.JWT())
	{
		admin_route.GET("/index", admin.AdminIndex)
	}
}
