package models

import (
	"time"
)

type DeployServerRelation struct {
	CreateTime time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	DeployId   int       `xorm:"not null pk comment('部署服务id') INT(11)"`
	ServerId   int       `xorm:"not null pk comment('服务器ID') INT(11)"`
}
