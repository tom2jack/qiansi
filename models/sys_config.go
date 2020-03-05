package models

import "time"

// SysConfig 内部回转配置
type SysConfig struct {
	Key        string    `gorm:"primary_key;column:key;type:varchar(50);not null"`
	Data       string    `gorm:"column:data;type:longtext"`       // 数据
	Readme     string    `gorm:"column:readme;type:varchar(255)"` // 字段说明
	CreateTime time.Time `gorm:"column:create_time;type:datetime"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp"`
}

// 根据Key获取配置
func (m *SysConfig) Get(key string) bool {
	m.Key = key
	db := Mysql.First(m)
	return db.Error == nil && db.RowsAffected > 0
}
