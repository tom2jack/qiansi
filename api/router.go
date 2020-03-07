package api

import (
	"gitee.com/zhimiao/qiansi/api/admin"
	"gitee.com/zhimiao/qiansi/api/client"
	_ "gitee.com/zhimiao/qiansi/docs"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

var Router *gin.Engine

func initRoute() {
	Router = gin.New()
	Router.Use(gin.Recovery(), logMiddleware())
	// 状态监控
	Router.Use(ginprom.PromMiddleware(nil))
	Router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	// 跨域支持
	Router.Use(corsMiddleware())

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
		// hook部署
		client_route.GET("/ApiDeployRun", client.ApiDeployRun)
		client_route.POST("/ApiDeployRun", client.ApiDeployRun)
		// 客户端交互请求
		client_route.Use(clientAuthMiddleware())
		{
			// 客户端部署日志推送
			client_route.POST("/ApiDeployLog", client.ApiDeployLog)
			// 获取任务
			client_route.GET("/ApiGetDeployTask", client.ApiGetDeployTask)
			// 部署回调
			client_route.GET("/ApiDeployNotify", client.ApiDeployNotify)
			// 获取监控配置
			client_route.GET("/ApiGetTelegrafConfig", client.ApiGetTelegrafConfig)
			// 客户端监控推送
			client_route.POST("/ApiClientMetric", client.ApiClientMetric)
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
		// 需要登陆的部分
		admin_route.Use(jwtMiddleware())
		{
			admin_route.POST("/UserResetPwd", admin.UserResetPwd)

			admin_route.GET("/ServerLists", admin.ServerLists)
			admin_route.POST("/ServerSet", admin.ServerSet)
			admin_route.DELETE("/ServerDel", admin.ServerDel)

			admin_route.GET("/DeployLists", admin.DeployLists)
			admin_route.DELETE("/DeployDel", admin.DeployDel)
			admin_route.POST("/DeploySet", admin.DeploySet)
			admin_route.POST("/DeployRelationServer", admin.DeployRelationServer)
			admin_route.GET("/DeployRunLogTab", admin.DeployRunLogTab)
			admin_route.GET("/DeployRunLog", admin.DeployRunLog)
			admin_route.GET("/DeployLog", admin.DeployLog)
			admin_route.POST("/DeployServer", admin.DeployServer)
			admin_route.GET("/DeployDo", admin.DeployDo)
			admin_route.GET("/DeployLink", admin.DeployLink)

			admin_route.GET("/ScheduleLists", admin.ScheduleLists)
			admin_route.POST("/ScheduleCreate", admin.ScheduleCreate)
			admin_route.DELETE("/ScheduleDel", admin.ScheduleDel)
			admin_route.GET("/ScheduleDo", admin.ScheduleDo)
		}
	}
}