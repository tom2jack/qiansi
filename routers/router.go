package routers

import (
	"github.com/gin-gonic/gin"
	"tools-server/controllers"
)

var (
	Router *gin.Engine = gin.Default()
)

func init() {

	Router.GET("/api/index", controllers.ApiIndex)
	Router.StaticFile("/", "assets/html/index.html")
	Router.Static("/static", "assets/html")
}
