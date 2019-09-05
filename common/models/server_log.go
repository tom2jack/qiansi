package models

import (
	"time"
)

type ServerLog struct {
	ClientIp   string    `xorm:"default '0.0.0.0' comment('客户端IP') VARCHAR(50)"`
	Content    string    `xorm:"comment('日志内容') TEXT"`
	CreateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' comment('创建时间') DATETIME"`
	DeviceId   string    `xorm:"not null default '' comment('设备号') CHAR(36)"`
	Id         int64     `xorm:"pk autoincr BIGINT(20)"`
	ServerId   int       `xorm:"not null default 0 comment('服务器编号') INT(11)"`
}
