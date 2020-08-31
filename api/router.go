package api

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/zhi-miao/qiansi/api/admin"
	"github.com/zhi-miao/qiansi/api/client"
	_ "github.com/zhi-miao/qiansi/docs"
)

var Router *gin.Engine

func initRoute() {
	Router = gin.New()
	Router.Use(gin.Recovery(), logMiddleware())
	// 状态监控
	ginpprof.Wrap(Router)
	// Router.Use(ginprom.PromMiddleware(nil))
	// Router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	// 跨域支持
	Router.Use(corsMiddleware())
	/* ------ 文档模块 ------- */
	Router.GET("/docs/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))
	/* ------ 静态页模块 ------- */
	Router.StaticFile("/", "assets/html/index.html")
	Router.Static("/static", "assets/html")
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
			adminRoute.GET("/DeployDetail", admin.Deploy.Detail)
			adminRoute.DELETE("/DeployDelete", admin.Deploy.Del)
			adminRoute.POST("/DeployCreate", admin.Deploy.Create)
			adminRoute.PUT("/DeployUpdate", admin.Deploy.Update)
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
