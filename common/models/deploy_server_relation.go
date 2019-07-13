package models

import (
	"time"
)

type DeployServerRelation struct {
	CreateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	DeployId      int       `xorm:"not null pk comment('部署服务id') INT(11)"`
	DeployVersion int       `xorm:"default 0 comment('已部署版本') INT(10)"`
	ServerId      int       `xorm:"not null pk comment('服务器ID') INT(11)"`
	UpdateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
