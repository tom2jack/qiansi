package errors

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MissingParameter int = iota
)

type ApiError interface {
	Error() string
	GetCode() int
}

type apiError struct {
	Msg         string
	Code        int
	MsgTemplate string
	Args        []interface{}
}

func (e *apiError) Error() string {
	return e.Msg
}

func (e apiError) GetCode() int {
	return e.Code
}

func NewApiError(code int, msg string, args ...interface{}) error {
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
