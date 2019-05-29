package main

import (
	"github.com/gin-gonic/gin"
	"tools-server/conf"
	"tools-server/routers"
)

func main() {
	conf.LoadConfig()
	gin.SetMode(conf.App.MustValue("server", "run_mode"))
	routers.Router.Run(conf.App.MustValue("server", "http_port", ":7091"))
}
