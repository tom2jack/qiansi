package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"qiansi/qiansi-server/api/admin"
	"qiansi/qiansi-server/api/client"
	_ "qiansi/qiansi-server/docs"
	"qiansi/qiansi-server/middleware"
)

var (
	Router *gin.Engine = gin.Default()
)

func LoadRouter() {
	/* ------ 文档模块 ------- */
	Router.GET("/docs/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	/* ------ 静态页模块 ------- */
	Router.StaticFile("/", "assets/html/index.html")
	Router.Static("/static", "assets/html")

	/* ------ 客户端模块 ------- */
	client_route := Router.Group("/client")
	{
		// 客户端注册
		client_route.GET("/ApiRegServer", client.ApiRegServer)
		// 客户端交互请求
		client_route.Use(middleware.ClientAuth())
		{
			client_route.GET("/ApiGetDeployTask", client.ApiGetDeployTask)
		}
	}

	/* ------ 后台模块 ------- */
	admin_route := Router.Group("/admin")
	{
		// 获取图片验证码
		admin_route.GET("/VerifyByImg", admin.VerifyByImg)
		// 获取短信验证码
		admin_route.POST("/VerifyBySMS", admin.VerifyBySMS)
		// 登录
		admin_route.POST("/UserSigin", admin.UserSigin)
		// 注册
		admin_route.POST("/UserSiginUp", admin.UserSiginUp)
		admin_route.GET("/index", admin.AdminIndex)
		// 需要登陆的部分
		admin_route.Use(middleware.JWT())
		{
			admin_route.POST("/UserResetPwd", admin.UserResetPwd)
			admin_route.GET("/ServerLists", admin.ServerLists)
			admin_route.DELETE("/ServerDel", admin.ServerDel)
			admin_route.GET("/DeployLists", admin.DeployLists)
			admin_route.DELETE("/DeployDel", admin.DeployDel)
			admin_route.POST("/DeploySet", admin.DeploySet)
			admin_route.POST("/DeployRelationServer", admin.DeployRelationServer)
			admin_route.DELETE("/DeployUnRelationServer", admin.DeployUnRelationServer)
			admin_route.GET("/DeployDo", admin.DeployDo)
		}
	}

}
