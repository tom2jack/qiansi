# 千丝平台

#### 介绍

千丝平台，集成代码发布、计划任务、配置管理等功能

> gorm数据库反转生成

```
go get github.com/xxjwxc/gormt

gormt -H "127.0.0.1" -u "qiansi_user" -p "root" -d "qiansi_dev" -o "./tools/out/models" --port=3306 -l json

```

> 文档自动生成

```
go get -u github.com/swaggo/swag/cmd/swag

swag init
```

> Redis字典

| key | 备注 |
|:------|:-------|
| QIANSI:verify:${string} | 验证码 |
| QIANSI:verify: phone:${phoneNumber} | 手机验证码 |
| QIANSI:schedule:server | 服务端计划任务 |

> 服务模式运行

`/usr/lib/systemd/system/qiansi-server.service`

```shell script
[Unit]
Description=The QIANSI Service
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
WorkingDirectory=/zhimiao/qiansi
PrivateTmp=true
Restart=always
Type=simple
ExecStart=/zhimiao/qiansi/qiansi_server_linux_amd64
ExecStop=/usr/bin/kill -15  $MAINPID

[Install]
WantedBy=multi-user.target
```
