package resp

import (
	"github.com/zhi-miao/qiansi/common/errors"
)

// PageInfo 分页返回标准结构
type PageInfo struct {
	Page      int         `json:"Page"`
	PageSize  int         `json:"PageSize"`
	TotalSize int         `json:"TotalSize"`
	Rows      interface{} `json:"Rows"`
}

// ApiErrorResult api默认错误返回结构
type ApiErrorResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Error(e error) *ApiErrorResult {
	result := &ApiErrorResult{}
	switch e := e.(type) {
	case errors.ApiError:
		result.Code = e.GetCode()
		result.Msg = e.Error()
	default:
		result.Code = 0
		result.Msg = e.Error()
	}
	return result
}
