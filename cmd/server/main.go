package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"tools-server/common/aliyun"
	"tools-server/conf"
	"tools-server/models"
	"tools-server/server/middleware"
	"tools-server/server/routers"
	"tools-server/server/service"
)

var g errgroup.Group

func init() {
	//加载配置
	conf.S = conf.LoadConfig("assets/config/server.ini")
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

// @title 纸喵 Tools-Server API
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
	https_listen := conf.S.MustValue("server", "https_listen")
	readTimeout := time.Duration(conf.S.MustInt64("server", "read_timeout", 60)) * time.Second
	writeTimeout := time.Duration(conf.S.MustInt64("server", "write_timeout", 60)) * time.Second
	http_server := &http.Server{
		Addr:           http_listen,
		Handler:        routers.Router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	https_server := &http.Server{
		Addr:           https_listen,
		Handler:        routers.Router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Start HTTP Service Listening %s", http_listen)
	r := gin.New()
	r.Use(middleware.TLS())
	r.GET("/*any", func(c *gin.Context) {
		c.String(301, "为了通信安全，请使用本站的https服务， http --> https")
	})
	g.Go(func() error {
		return http_server.ListenAndServe()
	})

	log.Printf("Start HTTPS Service Listening %s", https_listen)
	g.Go(func() error {
		return https_server.ListenAndServeTLS(
			conf.S.MustValue("server", "ssl_public_file"),
			conf.S.MustValue("server", "ssl_private_file"),
		)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Fatal()
}

func httpRedirecting() {
	https_listen := conf.S.MustValue("server", "https_listen")
	http_listen := conf.S.MustValue("server", "http_listen")
	secureMiddleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     https_listen,
	})
	var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("为了通信安全，请使用本站的https服务， http --> https"))
	})
	app := secureMiddleware.Handler(myHandler)
	// HTTP
	go func() {
		log.Fatal(http.ListenAndServe(http_listen, app))
	}()
}

func setLoger() {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	//gin.DisableConsoleColor()
	log_file := conf.S.MustValue("server", "log_file", "")
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
