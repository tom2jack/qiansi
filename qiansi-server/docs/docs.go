// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-07-13 18:14:27.5675468 +0800 CST m=+0.088001401

package docs

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "纸喵软件系列之服务端",
        "title": "纸喵 qiansi API",
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
        "/admin/DeployDel": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "删除部署服务",
                "parameters": [
                    {
                        "type": "string",
                        "description": "服务器ID",
                        "name": "deploy_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"操作成功\",\"data\": null}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/DeployDo": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "启动部署",
                "parameters": [
                    {
                        "type": "string",
                        "description": "部署应用ID",
                        "name": "deploy_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"启动成功\",\"data\": null}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/DeployLists": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取部署服务列表",
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"读取成功\",\"data\": [{\"AfterCommand\": \"324545\",\"BeforeCommand\": \"1232132132\",\"Branch\": \"123213\",\"CreateTime\": \"2019-02-28T10:24:41+08:00\",\"DeployType\": 1,\"Id\": 491,\"LocalPath\": \"123213\",\"NowVersion\": 0,\"RemoteUrl\": \"123213\",\"Title\": \"491-一号机器的修改241\",\"Uid\": 2,\"UpdateTime\": \"2019-02-28T10:25:17+08:00\"}]}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/DeployRelationServer": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "部署应用关联服务器",
                "parameters": [
                    {
                        "type": "string",
                        "description": "部署应用ID",
                        "name": "deploy_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "服务器ID",
                        "name": "server_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"关联成功\",\"data\": null}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/DeploySet": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "设置部署应用",
                "parameters": [
                    {
                        "type": "string",
                        "description": "应用ID",
                        "name": "Id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "后置命令",
                        "name": "AfterCommand",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "前置命令",
                        "name": "BeforeCommand",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "抓取分支",
                        "name": "Branch",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "部署方式",
                        "name": "DeployType",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "部署地址",
                        "name": "LocalPath",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "资源地址",
                        "name": "RemoteUrl",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "应用名称",
                        "name": "Title",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"操作成功\",\"data\": null}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/DeployUnRelationServer": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "部署应用取消关联服务器",
                "parameters": [
                    {
                        "type": "string",
                        "description": "部署应用ID",
                        "name": "deploy_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "服务器ID",
                        "name": "server_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"关联解除成功\",\"data\": null}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/ServerDel": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "删除服务器",
                "parameters": [
                    {
                        "type": "string",
                        "description": "服务器ID",
                        "name": "server_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"操作成功\",\"data\": null}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/ServerLists": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取服务器(客户端)列表",
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"读取成功\",\"data\": [{\"CreateTime\": \"2019-03-02T16:10:10+08:00\",\"DeviceId\": \"\",\"Domain\": \"127.0.0.1\",\"Id\": 1,\"ServerName\": \"纸喵5号机\",\"ServerRuleId\": 0,\"ServerStatus\": 0,\"Uid\": 2,\"UpdateTime\": \"2019-05-22T17:40:18+08:00\"}]}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/UserResetPwd": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "修改密码",
                "parameters": [
                    {
                        "type": "string",
                        "description": "旧密码",
                        "name": "old_password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "新密码",
                        "name": "new_password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"修改成功\", \"data\": null}}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/UserSigin": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "登录",
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
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"登录成功\", \"data\": {\"CreateTime\": \"2019-02-27T16:11:27+08:00\",\"InviterUid\": 0,\"Password\": \"\",\"Phone\": \"15061370322\",\"Status\": 1,\"Uid\": 2, \"UpdateTime\": \"2019-02-27T16:19:54+08:00\", \"Token\":\"sdfsdafsd..\"}}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/UserSiginUp": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "注册账号",
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
                        "description": "密码",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "短信验证码",
                        "name": "code",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "邀请人UID",
                        "name": "inviter_uid",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"注册成功\",\"data\": {\"CreateTime\": \"2019-02-27T16:11:27+08:00\",\"InviterUid\": 0,\"Password\": \"\",\"Phone\": \"15061370322\",\"Status\": 1,\"Uid\": 2, \"UpdateTime\": \"2019-02-27T16:19:54+08:00\", \"Token\":\"sdfsdafsd..\"}}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/VerifyByImg": {
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
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/admin/VerifyBySMS": {
            "post": {
                "consumes": [
                    "multipart/form-data"
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
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/clinet/ApiGetDeployTask": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取服务器部署任务清单",
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"读取成功\",\"data\": [deploy]}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        },
        "/clinet/ApiRegServer": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "服务器注册",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户UID",
                        "name": "uid",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客户端设备号",
                        "name": "device",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"登录成功\", \"data\": {\"CreateTime\": \"2019-02-27T16:11:27+08:00\",\"InviterUid\": 0,\"Password\": \"\",\"Phone\": \"15061370322\",\"Status\": 1,\"Uid\": 2, \"UpdateTime\": \"2019-02-27T16:19:54+08:00\", \"Token\":\"sdfsdafsd..\"}}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.ApiResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ApiResult": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{Schemes: []string{}}

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
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
