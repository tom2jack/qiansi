package models

import (
	"github.com/jinzhu/gorm"
)

// ListByDeployId 根据应用获取关联数据
func (m *DeployServerRelation) ListByDeployId() []DeployServerRelation {
	data := &[]DeployServerRelation{}
	Mysql.Where("deploy_id = ?", m.DeployID).Find(data)
	return *data
}

// ListByServerId 根据服务器获取关联数据
func (m *DeployServerRelation) ListByServerId() []DeployServerRelation {
	data := &[]DeployServerRelation{}
	Mysql.Where("server_id = ?", m.ServerID).Find(data)
	return *data
}

// Relation 关联服务器
func (m *DeployServerRelation) Relation(relation bool) bool {
	var db *gorm.DB
	if relation {
		db = Mysql.Save(m)
	} else {
		db = Mysql.Delete(m, "server_id=? and deploy_id=?", m.ServerID, m.DeployID)
	}
	return db.Error == nil && db.RowsAffected == 1
}
