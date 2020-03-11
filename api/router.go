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
		client_route.GET("/ApiRegServer", client.Api.RegServer)
		// hook部署
		client_route.GET("/ApiDeployRun", client.Api.DeployRun)
		client_route.POST("/ApiDeployRun", client.Api.DeployRun)
		// 客户端交互请求
		client_route.Use(clientAuthMiddleware())
		{
			// 客户端部署日志推送
			client_route.POST("/ApiDeployLog", client.Api.DeployLog)
			// 获取任务
			client_route.GET("/ApiGetDeployTask", client.Api.GetDeployTask)
			// 部署回调
			client_route.GET("/ApiDeployNotify", client.Api.DeployNotify)
			// 获取监控配置
			client_route.GET("/ApiGetTelegrafConfig", client.Api.GetTelegrafConfig)
			// 客户端监控推送
			client_route.POST("/ApiClientMetric", client.Api.ClientMetric)
		}
	}

	/* ------ 后台模块 ------- */
	admin_route := Router.Group("/admin")
	{
		// 获取图片验证码
		admin_route.GET("/VerifyByImg", admin.Verify.ByImg)
		// 获取短信验证码
		admin_route.POST("/VerifyBySMS", admin.Verify.BySMS)
		// 登录
		admin_route.POST("/UserSigin", admin.User.Sigin)
		// 注册
		admin_route.POST("/UserSiginUp", admin.User.SiginUp)
		// 需要登陆的部分
		admin_route.Use(jwtMiddleware())
		{
			admin_route.POST("/UserResetPwd", admin.User.ResetPwd)

			admin_route.GET("/ServerLists", admin.Server.Lists)
			admin_route.POST("/ServerSet", admin.Server.Set)
			admin_route.DELETE("/ServerDel", admin.Server.Del)

			admin_route.GET("/DeployLists", admin.Deploy.Lists)
			admin_route.DELETE("/DeployDel", admin.Deploy.Del)
			admin_route.POST("/DeploySet", admin.Deploy.Set)
			admin_route.POST("/DeployRelationServer", admin.Deploy.RelationServer)
			admin_route.GET("/DeployRunLogTab", admin.Deploy.RunLogTab)
			admin_route.GET("/DeployRunLog", admin.Deploy.RunLog)
			admin_route.GET("/DeployLog", admin.Deploy.Log)
			admin_route.POST("/DeployServer", admin.Deploy.Server)
			admin_route.GET("/DeployDo", admin.Deploy.Do)
			admin_route.GET("/DeployLink", admin.Deploy.Link)

			admin_route.GET("/ScheduleLists", admin.Schedule.Lists)
			admin_route.POST("/ScheduleCreate", admin.Schedule.Create)
			admin_route.DELETE("/ScheduleDel", admin.Schedule.Del)
			admin_route.GET("/ScheduleDo", admin.Schedule.Do)
		}
	}
}
