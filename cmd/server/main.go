package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"qiansi/common/aliyun"
	"qiansi/common/conf"
	"qiansi/common/models"
	"qiansi/common/zmlog"
	"qiansi/qiansi-server/routers"
	"qiansi/qiansi-server/service"
	"time"
)

var g errgroup.Group

func init() {
	//加载配置
	conf.S = conf.LoadConfig("assets/config/server.ini")
	// 配置日志记录方式
	zmlog.InitLog("server.log")
	//加载路由
	routers.LoadRouter()
	//启动服务
	service.LoadService()
	//初始化Redis
	models.LoadRedis()
	//初始化MySQL
	models.LoadMysql()
	// 加载阿里云SDK
	aliyun.LoadAliyunSDK()
}

// @title 纸喵 qiansi API
// @version 1.0
// @description 纸喵软件系列之服务端
// @termsOfService http://zhimiao.org

// @contact.name API Support
// @contact.url http://tools.zhimiao.org
// @contact.email mail@xiaoliu.org

// @host localhost:8000
// @BasePath
func main() {
	defer destroy()
	gin.SetMode(conf.S.MustValue("server", "run_mode"))
	http_listen := conf.S.MustValue("server", "http_listen")
	readTimeout := time.Duration(conf.S.MustInt64("server", "read_timeout", 60)) * time.Second
	writeTimeout := time.Duration(conf.S.MustInt64("server", "write_timeout", 60)) * time.Second
	http_server := &http.Server{
		Addr:           http_listen,
		Handler:        routers.Router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	zmlog.Info("Start HTTP Service Listening %s", http_listen)
	http_server.ListenAndServe()
}

//Destroy 销毁资源
func destroy() {
	models.ZM_Mysql.Close()
}
