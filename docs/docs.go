// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-06-11 09:10:33.1978751 +0800 CST m=+0.087998101

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "纸喵软件系列之服务端",
        "title": "纸喵 Tools-Server API",
        "termsOfService": "http://zhimiao.org",
        "contact": {
            "name": "API Support",
            "url": "http://tools.zhimiao.org",
            "email": "mail@xiaoliu.org"
        },
        "license": {},
        "version": "1.0"
    },
    "host": "localhost:8000",
    "paths": {
        "/admin/verify/VerifyByImg": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取图片验证码",
                "responses": {
                    "200": {
                        "description": "{\"code\":1,\"msg\":\"\",\"data\":{\"idkey\":\"ckFbFAcMo7sy7qGyonAd\",\"img\":\"data:image/png;base64,iVBORw0...\"}}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/verify/VerifyBySMS": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取短信验证码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "手机号",
                        "name": "phone",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "图片验证码句柄",
                        "name": "img_idkey",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "图片验证码",
                        "name": "img_code",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":1,\"msg\":\"发送成功\",\"data\":null}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
