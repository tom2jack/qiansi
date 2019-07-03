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

func logDo(stat string, format string, v ...interface{}) {
	format = fmt.Sprintf("\n%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), stat, format)
	if stat == "ERROR" {
		fmt.Fprintf(gin.DefaultErrorWriter, format, v...)
		os.Exit(0)
	}
	fmt.Fprintf(gin.DefaultErrorWriter, format, v...)
}

func Error(format string, v ...interface{}) {
	logDo("ERROR", format, v...)
}

func Warn(format string, v ...interface{}) {
	logDo("WARN", format, v...)
}

func Info(format string, v ...interface{}) {
	logDo("INFO", format, v...)
}

func Debug(format string, v ...interface{}) {
	logDo("DEBUG", format, v...)
}

func Trace(format string, v ...interface{}) {
	logDo("TRACE", format, v...)
}
