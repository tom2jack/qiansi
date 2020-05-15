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
	clientRoute := Router.Group("/client")
	{
		clientRoute.GET("/ApiRegServer", client.Api.RegServer)
		clientRoute.GET("/ApiDeployRun", client.Api.DeployRun)
		clientRoute.POST("/ApiDeployRun", client.Api.DeployRun)
		// 客户端交互请求
		clientRoute.Use(clientAuthMiddleware())
		{
			clientRoute.POST("/ApiDeployLog", client.Api.DeployLog)
			clientRoute.GET("/ApiGetDeployTask", client.Api.GetDeployTask)
			clientRoute.GET("/ApiDeployNotify", client.Api.DeployNotify)
			clientRoute.GET("/ApiGetTelegrafConfig", client.Api.GetTelegrafConfig)
			clientRoute.POST("/ApiClientMetric", client.Api.ClientMetric)
		}
	}
	/* ------ 后台模块 ------- */
	adminRoute := Router.Group("/admin")
	{
		adminRoute.GET("/VerifyByImg", admin.Verify.ByImg)
		adminRoute.POST("/VerifyBySMS", admin.Verify.BySMS)
		adminRoute.POST("/UserSigin", admin.User.Sigin)
		adminRoute.POST("/UserSiginUp", admin.User.SiginUp)
		adminRoute.Use(jwtMiddleware())
		{
			adminRoute.POST("/UserResetPwd", admin.User.ResetPwd)

			adminRoute.GET("/ServerLists", admin.Server.Lists)
			adminRoute.POST("/ServerSet", admin.Server.Set)
			adminRoute.DELETE("/ServerDel", admin.Server.Del)

			adminRoute.GET("/DeployLists", admin.Deploy.Lists)
			adminRoute.DELETE("/Deploy", admin.Deploy.Del)
			adminRoute.POST("/Deploy", admin.Deploy.Create)
			adminRoute.PUT("/Deploy", admin.Deploy.Update)
			adminRoute.GET("/DeployRunLogTab", admin.Deploy.RunLogTab)
			adminRoute.GET("/DeployRunLog", admin.Deploy.RunLog)
			adminRoute.GET("/DeployLog", admin.Deploy.Log)
			adminRoute.GET("/DeployServer", admin.Deploy.DeployServer)
			adminRoute.GET("/DeployDo", admin.Deploy.Do)
			adminRoute.GET("/DeployLink", admin.Deploy.Link)

			adminRoute.GET("/ScheduleLists", admin.Schedule.Lists)
			adminRoute.POST("/ScheduleCreate", admin.Schedule.Create)
			adminRoute.DELETE("/ScheduleDel", admin.Schedule.Del)
			adminRoute.GET("/ScheduleDo", admin.Schedule.Do)

			adminRoute.GET("/dashboard/info", admin.Dashboard.Info)
			adminRoute.GET("/dashboard/IndexMetric", admin.Dashboard.IndexMetric)
		}
	}
}
