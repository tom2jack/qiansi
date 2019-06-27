package models

import (
	"time"
)

type Server struct {
	ApiSecret    string    `xorm:"not null default '' comment('API密钥') VARCHAR(32)"`
	CreateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	DeviceId     string    `xorm:"not null comment('服务器唯一设备号') CHAR(36)"`
	Domain       string    `xorm:"not null comment('服务器地址(域名/ip)') VARCHAR(255)"`
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	ServerName   string    `xorm:"not null default '' comment('服务器备注名') VARCHAR(64)"`
	ServerRuleId int       `xorm:"not null default 0 comment('服务器规则id') INT(11)"`
	ServerStatus int       `xorm:"not null default 0 comment('服务器状态 -1-失效 0-待认领 1-已分配通信密钥 2-已绑定') TINYINT(1)"`
	Uid          int       `xorm:"not null comment('用户ID') index INT(10)"`
	UpdateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
