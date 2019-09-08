package utils

import (
	"fmt"
	"gopkg.in/go-playground/validator.v8"
)

func Validator(err error) string {
	if v8, ok := err.(validator.ValidationErrors); ok {
		for _, v := range v8 {
			return fmt.Sprintf("%s参数%s规则校验失败", v.Field, v.ActualTag)
		}
	}
	return "参数绑定错误"
}
