package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qiansi/common/conf"
	"qiansi/common/logger"
	"qiansi/qiansi-server/api"
	_ "qiansi/qiansi-server/schedule"
	_ "qiansi/qiansi-server/udp_service"
	"time"
)

// @title 纸喵 qiansi API
// @version 1.0
// @description 纸喵软件系列之服务端
// @termsOfService http://zhimiao.org

// @contact.name API Support
// @contact.url http://tools.zhimiao.org
// @contact.email mail@xiaoliu.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1315
// @basepath /

func main() {
	gin.SetMode(conf.S.MustValue("server", "run_mode"))
	httpListen := conf.S.MustValue("server", "http_listen")
	httpServer := &http.Server{
		Addr:           httpListen,
		Handler:        api.Router,
		ReadTimeout:    time.Duration(conf.S.MustInt64("server", "read_timeout", 60)) * time.Second,
		WriteTimeout:   time.Duration(conf.S.MustInt64("server", "write_timeout", 60)) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Info("Start HTTP Service Listening %s", httpListen)
	httpServer.ListenAndServe()
}
