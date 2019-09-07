package models

import (
	"time"
)

type Member struct {
	CreateTime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	Id         int       `xorm:"not null pk autoincr INT(10)"`
	InviterUid int       `xorm:"not null default 0 comment('邀请人UID') INT(10)"`
	Password   string    `xorm:"not null comment('密码') VARCHAR(255)"`
	Phone      string    `xorm:"not null comment('手机号') unique CHAR(11)"`
	Status     int       `xorm:"not null default 1 comment('0-锁定 1-正常') TINYINT(1)"`
	UpdateTime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
}
