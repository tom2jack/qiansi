package models

import "time"

// Telegraf 监控配置表
type Telegraf struct {
	ServerID   int       `gorm:"primary_key;column:server_id;type:int(11);not null"` // 服务器ID
	UId        int       `gorm:"column:uid;type:int(11)"`                            // 用户ID
	TomlConfig string    `gorm:"column:toml_config;type:longtext"`                   // 私有配置
	IsOpen     int8      `gorm:"column:is_open;type:tinyint(1)"`                     // 是否开启监控功能 1-开启 2-关闭
	CreateTime time.Time `gorm:"column:create_time;type:datetime"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime"`
}

// 根据Key获取配置
func (m *Telegraf) Get() bool {
	db := Mysql.Model(m)
	if m.UId > 0 && m.ServerID > 0 {
		db = db.Select("uid=?", m.UId).First(m)
		return db.Error == nil && db.RowsAffected > 0
	}
	return false
}
