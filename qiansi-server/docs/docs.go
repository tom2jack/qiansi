// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-08-17 19:00:45.9928614 +0800 CST m=+0.524206201

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
                "summary": "删除部署应用",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.DeployDelParam"
                        }
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
                "summary": "启动部署 TODO: 后期关闭此接口的开放特性，新增外部接口，通过不可枚举key作为部署参数",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.DeployDoParam"
                        }
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
                "summary": "获取部署应用列表",
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
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.DeployRelationParam"
                        }
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
        "/admin/DeployServer": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取当前部署应用的服务器列表",
                "responses": {
                    "200": {
                        "description": "返回",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.Server"
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
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.DeploySetParam"
                        }
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
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.ServerDelParam"
                        }
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
            "get": {
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
                            "$ref": "#/definitions/models.Server"
                        }
                    }
                }
            }
        },
        "/admin/ServerSet": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "设置服务器信息",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.DeploySetParam"
                        }
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
        "/admin/UserResetPwd": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "修改密码",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.UserResetPwdParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"修改成功\", \"data\": null}",
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
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.UserSiginParam"
                        }
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
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "注册账号",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.UserSiginUpParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 1,\"msg\": \"注册成功\",\"data\": {\"CreateTime\": \"2019-02-27T16:11:27+08:00\",\"InviterUid\": 0,\"Password\": \"\",\"Phone\": \"15061370322\",\"Status\": 1,\"Uid\": 2, \"UpdateTime\": \"2019-02-27T16:19:54+08:00\", \"Token\":\"sdfsdafsd..\"}}",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zresp.UserInfoVO"
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
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "获取短信验证码",
                "parameters": [
                    {
                        "description": "入参集合",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/zreq.VerifyBySMSParam"
                        }
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
        "/clinet/ApiDeployNotify": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "客户端部署成功回调",
                "parameters": [
                    {
                        "type": "string",
                        "description": "版本号",
                        "name": "version",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "部署应用ID",
                        "name": "deploy_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
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
            "get": {
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
        },
        "/clinet/LogPush": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "客户端日志推送",
                "parameters": [
                    {
                        "type": "string",
                        "description": "客户端平台编号",
                        "name": "server_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客户端设备号",
                        "name": "device",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "日志文本内容",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
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
        },
        "models.Server": {
            "type": "object",
            "properties": {
                "apiSecret": {
                    "type": "string"
                },
                "createTime": {
                    "type": "string"
                },
                "deviceId": {
                    "type": "string"
                },
                "domain": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "serverName": {
                    "type": "string"
                },
                "serverRuleId": {
                    "type": "integer"
                },
                "serverStatus": {
                    "type": "integer"
                },
                "uid": {
                    "type": "integer"
                },
                "updateTime": {
                    "type": "string"
                }
            }
        },
        "zreq.DeployDelParam": {
            "type": "object",
            "properties": {
                "deployId": {
                    "description": "部署应用ID",
                    "type": "integer"
                }
            }
        },
        "zreq.DeployDoParam": {
            "type": "object",
            "properties": {
                "deployId": {
                    "description": "部署应用ID",
                    "type": "integer"
                }
            }
        },
        "zreq.DeployRelationParam": {
            "type": "object",
            "properties": {
                "deployId": {
                    "description": "部署应用ID",
                    "type": "integer"
                },
                "relation": {
                    "description": "关联操作 true-关联 false-取消",
                    "type": "boolean"
                },
                "serverId": {
                    "description": "服务器ID",
                    "type": "integer"
                }
            }
        },
        "zreq.DeploySetParam": {
            "type": "object",
            "properties": {
                "afterCommand": {
                    "description": "后置命令",
                    "type": "string"
                },
                "beforeCommand": {
                    "description": "前置命令",
                    "type": "string"
                },
                "branch": {
                    "description": "git分支",
                    "type": "string"
                },
                "deployKeys": {
                    "description": "部署私钥",
                    "type": "string"
                },
                "deployType": {
                    "description": "部署类型 0-本地 1-git 2-zip",
                    "type": "integer"
                },
                "id": {
                    "description": "应用唯一编号",
                    "type": "integer"
                },
                "localPath": {
                    "description": "本地部署地址",
                    "type": "string"
                },
                "remoteUrl": {
                    "description": "资源地址",
                    "type": "string"
                },
                "title": {
                    "description": "应用名称",
                    "type": "string"
                }
            }
        },
        "zreq.ServerDelParam": {
            "type": "object",
            "properties": {
                "serverId": {
                    "type": "integer"
                }
            }
        },
        "zreq.UserResetPwdParam": {
            "type": "object",
            "properties": {
                "newPassword": {
                    "type": "string"
                },
                "oldPassword": {
                    "type": "string"
                }
            }
        },
        "zreq.UserSiginParam": {
            "type": "object",
            "properties": {
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "phone": {
                    "description": "手机号",
                    "type": "string"
                }
            }
        },
        "zreq.UserSiginUpParam": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "短信验证码",
                    "type": "string"
                },
                "inviterUid": {
                    "description": "邀请人",
                    "type": "integer"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "phone": {
                    "description": "手机号",
                    "type": "string"
                }
            }
        },
        "zreq.VerifyBySMSParam": {
            "type": "object",
            "properties": {
                "imgCode": {
                    "description": "验证码code",
                    "type": "string"
                },
                "imgIdKey": {
                    "description": "图片验证码id",
                    "type": "string"
                },
                "phone": {
                    "description": "手机号",
                    "type": "string"
                }
            }
        },
        "zresp.UserInfoVO": {
            "type": "object",
            "properties": {
                "token": {
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
