// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-06-06 17:00:14.2737334 +0800 CST m=+0.101000101

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
    "host": "127.0.0.1",
    "paths": {
        "/admin/verify/VerifyByImg": {
            "get": {
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
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "图片验证码句柄",
                        "name": "img_idkey",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "图片验证码",
                        "name": "img_code",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":1,",
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
