package models

import (
	"time"
)

type DeployLog struct {
	Id            int       `xorm:"not null pk autoincr INT(11)"`
	Uid           int       `xorm:"not null comment('用户ID') INT(11)"`
	DeployId      int       `xorm:"not null comment('服务器ID') INT(11)"`
	ServerId      int       `xorm:"not null comment('服务器ID') INT(11)"`
	DeviceId      string    `xorm:"not null comment('服务器唯一设备号') CHAR(36)"`
	DeployVersion int       `xorm:"not null comment('服务器ID') INT(11)"`
	Content       string    `xorm:"not null default '' comment('日志内容') TEXT"`
	ClientIp      string    `xorm:"not null comment('客户端IP') varchar(50)"`
	CreateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
