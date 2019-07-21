package models

import (
	"time"
)

type ServerLog struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	ServerId   int       `xorm:"not null comment('服务器ID') INT(11)"`
	DeviceId   string    `xorm:"not null comment('服务器唯一设备号') CHAR(36)"`
	Content    string    `xorm:"not null default '' comment('日志内容') TEXT"`
	ClientIp   string    `xorm:"not null comment('客户端IP') varchar(50)"`
	CreateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
