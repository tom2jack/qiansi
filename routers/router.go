package routers

import (
	"github.com/gin-gonic/gin"
	"tools-server/controllers"
)

var (
	Router *gin.Engine = gin.Default()
)

func init() {
	Router.GET("/", controllers.Index)
}
