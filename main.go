package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"tools-server/conf"
	"tools-server/routers"
)

func main() {
	conf.LoadConfig()
	gin.SetMode(conf.App.MustValue("server", "run_mode"))
	endPoint := conf.App.MustValue("server", "http_port", ":7091")
	//readTimeout := conf.App.MustInt64("server", "read_timeout", 60)
	//writeTimeout := conf.App.MustInt64("server", "write_timeout", 60)

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routers.Router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
	routers.Router.Run()
}
