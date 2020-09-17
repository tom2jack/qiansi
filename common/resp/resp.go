package resp

import (
	"net/http"

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

func ApiSuccess(e interface{}) (int, interface{}) {
	return http.StatusOK, e
}

func ApiError(e ...interface{}) (int, *ApiErrorResult) {
	httpStatus := 400
	result := &ApiErrorResult{
		Code: 0,
		Msg:  "未知错误",
	}
	argLen := len(e)
	if argLen == 0 {
		return httpStatus, result
	}
	if argLen > 1 {
		if arg2, ok := e[1].(string); ok {
			result.Msg = arg2
		}
	}
	switch arg1 := e[0].(type) {
	case int:
		result.Code = arg1
	case errors.ApiErrorCode:
		result.Code = int(arg1)
	case string:
		result.Msg = arg1
	case errors.ApiError:
		result.Code = arg1.GetCode()
		result.Msg = arg1.Error()
	case error:
		result.Code = 0
		result.Msg = arg1.Error()
	}
	if s, ok := errors.ApiErrorHttpCodeRelation[errors.ApiErrorCode(result.Code)]; ok {
		httpStatus = s
	}
	return httpStatus, result
}
