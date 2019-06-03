package utils

import "github.com/gin-gonic/gin"

type Json struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func Show(C *gin.Context, code int, msg string, data interface{}) {
	C.JSON(200, Json{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	return
}
