package models

import "github.com/jinzhu/gorm"

type sysConfigModels struct {
	db *gorm.DB
}

func GetSysConfigModels() *sysConfigModels {
	return &sysConfigModels{
		db: Mysql,
	}
}

func (m *sysConfigModels) SetDB(db *gorm.DB) *sysConfigModels {
	m.db = db
	return m
}

// GetTelegraf Telegraf系统配置
func (m *sysConfigModels) GetTelegraf() string {
	cnf := &SysConfig{Key: "telegraf_config"}
	m.db.Model(&cnf).First(&cnf)
	return cnf.Data
}
