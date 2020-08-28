package resp

import (
	"github.com/gin-gonic/gin"
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
