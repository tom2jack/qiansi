package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"qiansi/models"
)

type ApiResponseJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func Show(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(200, ApiResponseJson{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	return
}

func ClientShow(c *gin.Context, code int, msg string, data interface{}) {
	json_str, err := json.Marshal(ApiResponseJson{
		Code: code,
		Msg:  msg,
		Data: data,
	})
	if err != nil {
		c.String(500, "ERROR!")
		return
	}
	server_id := c.GetInt("SERVER-ID")
	server := &models.Server{}
	models.ZM_Mysql.Select("api_secret").Limit(1).Find(server, server_id)
	result := EncyptogAES(string(json_str), server.ApiSecret)
	// result = base64.StdEncoding.EncodeToString([]byte(result))
	c.Header("ZHIMIAO-Encypt", "1")
	c.String(200, result)
}
