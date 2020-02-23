package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type DeployServerRelation struct {
	CreateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	DeployId      int       `xorm:"not null pk comment('部署服务id') INT(11)"`
	DeployVersion int       `xorm:"default 0 comment('已部署版本') INT(10)"`
	ServerId      int       `xorm:"not null pk comment('服务器ID') INT(11)"`
	UpdateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

// ListByDeployId 根据应用获取关联数据
func (m *DeployServerRelation) ListByDeployId() []DeployServerRelation {
	data := &[]DeployServerRelation{}
	Mysql.Where("deploy_id = ?", m.DeployId).Find(data)
	return *data
}

// ListByServerId 根据服务器获取关联数据
func (m *DeployServerRelation) ListByServerId() []DeployServerRelation {
	data := &[]DeployServerRelation{}
	Mysql.Where("server_id = ?", m.ServerId).Find(data)
	return *data
}

// Relation 关联服务器
func (m *DeployServerRelation) Relation(relation bool) bool {
	var db *gorm.DB
	if relation {
		db = Mysql.Save(&DeployServerRelation{
			ServerId: m.ServerId,
			DeployId: m.DeployId,
		})
	} else {
		db = Mysql.Delete(m, "server_id=? and deploy_id=?", m.ServerId, m.DeployId)
	}
	return db.Error == nil && db.RowsAffected == 1
}
