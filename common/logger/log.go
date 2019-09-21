package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

var logFile = "qiansi.log"

func init() {
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	//gin.DisableConsoleColor()
	_, err := os.Stat(logFile)
	var f *os.File
	if err != nil {
		// 记录到文件。
		f, _ = os.Create(logFile)
	} else {
		f, _ = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	}
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func logDo(stat string, format string, v ...interface{}) {
	format = fmt.Sprintf("\n%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), stat, format)
	fmt.Fprintf(gin.DefaultWriter, format, v...)
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
