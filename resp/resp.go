package resp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zhi-miao/qiansi/common/utils"
	"github.com/zhi-miao/qiansi/models"
)

// ApiResult api默认返回结构
type ApiResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// PageInfo 分页返回标准结构
type PageInfo struct {
	Page      int         `json:"Page"`
	PageSize  int         `json:"PageSize"`
	TotalSize int         `json:"TotalSize"`
	Rows      interface{} `json:"Rows"`
}

// NewApiResult 初始化api返回
func NewApiResult(arg ...interface{}) *ApiResult {
	result := &ApiResult{
		Code: 1,
		Msg:  "操作成功",
	}
	for k, v := range arg {
		if k == 0 {
			if v1, ok := v.(int); ok {
				result.setCode(v1)
			}
			if v1, ok := v.(error); ok {
				result.setCode(-5)
				result.setMsg(v1.Error())
			}
		}
		if k == 1 {
			if v2, ok := v.(string); ok {
				result.setMsg(v2)
			}
		}
		if k == 2 && v != nil {
			result.setData(v)
		}
	}
	return result
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

func (r *ApiResult) Json(c *gin.Context) {
	c.JSON(200, r)
}

func (r *ApiResult) Encypt(c *gin.Context) {
	json_str, err := json.Marshal(r)
	if err != nil {
		c.String(500, "ERROR!")
		return
	}
	server_id := c.GetInt("SERVER-ID")
	server := &models.Server{}
	result := utils.EncyptogAES(string(json_str), server.GetApiSecret(server_id))
	// result = base64.StdEncoding.EncodeToString([]byte(result))
	c.Header("ZHIMIAO-Encypt", "1")
	c.String(200, result)
}
