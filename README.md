# 千丝平台

#### 介绍

千丝平台，集成代码发布、计划任务、配置管理、服务器监控报警等功能

UI界面开源地址： [https://github.com/zhi-miao/qiansi-vue](https://github.com/zhi-miao/qiansi-vue)

客户端开源地址： [https://github.com/zhi-miao/qiansi-client](https://github.com/zhi-miao/qiansi-client)

官方cloud版本：[http://qiansi.zhimiao.org](http://qiansi.zhimiao.org)

官网介绍：[http://zhimiao.org/post/qiansi-news](http://zhimiao.org/post/qiansi-news)

#### 开发文档

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
