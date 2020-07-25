package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common"
	"net/http"
	"time"
)

func Start() {
	gin.SetMode(gin.DebugMode)
	// 初始化route
	initRoute()
	httpServer := &http.Server{
		Addr:           common.Config.Server.APIListen,
		Handler:        Router,
		ReadTimeout:    time.Duration(common.Config.Server.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(common.Config.Server.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Infof("Start HTTP Service Listening %s", common.Config.Server.APIListen)
	httpServer.ListenAndServe()
}
