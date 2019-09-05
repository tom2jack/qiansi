package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"qiansi/common/conf"
	"qiansi/common/zmlog"
	"qiansi/qiansi-server/routers"
	"time"
)

var g errgroup.Group

// @title 纸喵 qiansi API
// @version 1.0
// @description 纸喵软件系列之服务端
// @termsOfService http://zhimiao.org

// @contact.name API Support
// @contact.url http://tools.zhimiao.org
// @contact.email mail@xiaoliu.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:1315
// @basepath /

func main() {
	defer destroy()
	gin.SetMode(conf.S.MustValue("server", "run_mode"))
	httpListen := conf.S.MustValue("server", "http_listen")
	readTimeout := time.Duration(conf.S.MustInt64("server", "read_timeout", 60)) * time.Second
	writeTimeout := time.Duration(conf.S.MustInt64("server", "write_timeout", 60)) * time.Second
	httpServer := &http.Server{
		Addr:           httpListen,
		Handler:        routers.Router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	zmlog.Info("Start HTTP Service Listening %s", httpListen)
	httpServer.ListenAndServe()
}

//Destroy 销毁资源
func destroy() {

}
