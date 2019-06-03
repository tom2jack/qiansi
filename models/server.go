package models

import (
	"time"
)

type Server struct {
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	Uid          int       `xorm:"not null comment('用户ID') index INT(10)"`
	ServerName   string    `xorm:"not null default '' comment('服务器备注名') VARCHAR(64)"`
	ServerStatus int       `xorm:"not null default 0 comment('服务器状态 -1-失效 0-待认领 1-已分配通信密钥 2-已绑定') TINYINT(1)"`
	ApiSecret    string    `xorm:"not null default '' comment('API密钥') VARCHAR(255)"`
	ServerRuleId int       `xorm:"not null default 0 comment('服务器规则id') INT(11)"`
	DeviceId     string    `xorm:"not null comment('服务器唯一设备号') CHAR(36)"`
	Domain       string    `xorm:"not null comment('服务器地址(域名/ip)') VARCHAR(255)"`
	CreateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	UpdateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
