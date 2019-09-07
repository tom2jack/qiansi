package models

import (
	"time"
)

type DeployLog struct {
	ClientIp      string    `xorm:"default '0.0.0.0' comment('客户端IP') VARCHAR(50)"`
	Content       string    `xorm:"not null comment('日志内容') TEXT"`
	CreateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' comment('创建时间') DATETIME"`
	DeployId      int       `xorm:"not null default 0 comment('部署ID') index(IX_search) INT(10)"`
	DeployVersion int       `xorm:"not null default 0 comment('部署版本') index(IX_search) INT(10)"`
	DeviceId      string    `xorm:"not null default '' comment('设备号') CHAR(36)"`
	Id            int64     `xorm:"pk autoincr BIGINT(20)"`
	ServerId      int       `xorm:"not null default 0 comment('服务器ID') index(IX_search) INT(11)"`
	Uid           int       `xorm:"not null default 0 comment('用户ID') index(IX_search) INT(10)"`
}
