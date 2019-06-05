package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"tools-server/common/aliyun"
	"tools-server/conf"
	"tools-server/models"
	"tools-server/routers"
	"tools-server/service"
)

func init() {
	//加载配置
	conf.LoadConfig()
	// 配置日志记录方式
	setLoger()
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

func main() {
	defer destroy()
	gin.Logger()
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
	log.Printf("Start HTTP Service Listening %s", endPoint)
	server.ListenAndServe()
}

func setLoger() {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	//gin.DisableConsoleColor()
	log_file := conf.App.MustValue("server", "log_file", "")
	if log_file != "" {
		// 记录到文件。
		f, _ := os.Create("gin.log")
		// 如果需要同时将日志写入文件和控制台，请使用以下代码。
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
}

//Destroy 销毁资源
func destroy() {
	models.ZM_Mysql.Close()
}
