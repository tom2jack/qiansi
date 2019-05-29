package main

import (
	"github.com/gin-gonic/gin"
	"tools-server/routers"
)

func main() {
	gin.SetMode(gin.DebugMode)
	routers.Router.Run(":7091")
}
