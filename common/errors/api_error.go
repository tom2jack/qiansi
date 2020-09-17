package errors

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	// 参数缺失
	MissingParameter ApiErrorCode = iota + 1
	// 最大限制
	MaximumLimit
	// 最小限制
	MinimumLimit
	// 参数格式错误
	InvalidArguments
	// 服务错误
	Service
	// 模型
	Models
	// 服务器异常
	InternalServerError
	// 未授权
	Unauthorized
)

var ApiErrorHttpCodeRelation = map[ApiErrorCode]int{
	MissingParameter:    http.StatusBadRequest,
	InternalServerError: http.StatusInternalServerError,
	Unauthorized:        http.StatusUnauthorized,
}

type ApiErrorCode int

type ApiError interface {
	Error() string
	GetCode() int
}

type apiError struct {
	Msg         string
	Code        ApiErrorCode
	MsgTemplate string
	Args        []interface{}
}

func (e *apiError) Error() string {
	return e.Msg
}

func (e apiError) GetCode() ApiErrorCode {
	return e.Code
}

func NewApiError(code ApiErrorCode, msg string, args ...interface{}) error {
	e := &apiError{
		Code: code,
		Msg:  msg,
		Args: args,
	}
	if len(args) > 0 {
		e.MsgTemplate = msg
		for i, v := range args {
			e.Msg = strings.ReplaceAll(e.Msg, "{"+strconv.Itoa(i)+"}", fmt.Sprint(v))
		}
	}
	return e
}
