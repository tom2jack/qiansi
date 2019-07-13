package models

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"qiansi/common/utils"
)

type ApiResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewApiResult(arg ...interface{}) *ApiResult {
	return &ApiResult{
		Code: 1,
		Msg:  "操作成功",
	}
}

func (r *ApiResult) setData(data interface{}) *ApiResult {
	r.Data = data
	return r
}

func (r *ApiResult) setMsg(msg string) *ApiResult {
	r.Msg = msg
	return r
}

func (r *ApiResult) setCode(code int) *ApiResult {
	r.Code = code
	return r
}

func (r *ApiResult) Success(c *gin.Context) {
	c.JSON(200, r)
}

func (r *ApiResult) Client(c *gin.Context) {
	json_str, err := json.Marshal(r)
	if err != nil {
		c.String(500, "ERROR!")
		return
	}
	server_id := c.GetInt("SERVER-ID")
	server := &Server{}
	ZM_Mysql.Select("api_secret").Limit(1).Find(server, server_id)
	result := utils.EncyptogAES(string(json_str), server.ApiSecret)
	// result = base64.StdEncoding.EncodeToString([]byte(result))
	c.Header("ZHIMIAO-Encypt", "1")
	c.String(200, result)
}
