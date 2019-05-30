package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"tools-server/conf"
	"tools-server/models"
	"tools-server/routers"
)

func init() {
	conf.LoadConfig()
	models.LoadRedis()
	models.LoadMysql()
}

func main() {
	gin.SetMode(conf.App.MustValue("server", "run_mode"))
	endPoint := conf.App.MustValue("server", "http_port", ":7091")
	readTimeout := time.Duration(conf.App.MustInt64("server", "read_timeout", 60)) * time.Second
	writeTimeout := time.Duration(conf.App.MustInt64("server", "write_timeout", 60)) * time.Second

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routers.Router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[info] start http server listening %s", endPoint)
	server.ListenAndServe()
}
