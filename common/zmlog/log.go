package zmlog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

func InitLog(logfile string) {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	//gin.DisableConsoleColor()
	if logfile != "" {
		// 记录到文件。
		f, _ := os.Create(logfile)
		// 如果需要同时将日志写入文件和控制台，请使用以下代码。
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
}

func Error(format string, v ...interface{}) {
	format = time.Now().Format("2006-01-02 15:04:05") + " [ERROR] " + format + "\n"
	fmt.Fprintf(gin.DefaultErrorWriter, format, v)
	os.Exit(0)
}

func Warn(format string, v ...interface{}) {
	format = time.Now().Format("2006-01-02 15:04:05") + " [WARN] " + format + "\n"
	fmt.Fprintf(gin.DefaultErrorWriter, format, v)
}

func Info(format string, v ...interface{}) {
	format = time.Now().Format("2006-01-02 15:04:05") + " [INFO] " + format + "\n"
	fmt.Fprintf(gin.DefaultErrorWriter, format, v)
}

func Debug(format string, v ...interface{}) {
	format = time.Now().Format("2006-01-02 15:04:05") + " [DEBUG] " + format + "\n"
	fmt.Fprintf(gin.DefaultErrorWriter, format, v)
}

func Trace(format string, v ...interface{}) {
	format = time.Now().Format("2006-01-02 15:04:05") + " [TRACE] " + format + "\n"
	fmt.Fprintf(gin.DefaultErrorWriter, format, v)
}
