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

func (m *sysConfigModels) get(key string) string {
	cnf := &SysConfig{Key: key}
	m.db.Model(&cnf).First(&cnf)
	return cnf.Data
}

// GetTelegraf Telegraf系统配置
func (m *sysConfigModels) GetTelegraf() string {
	return m.get("telegraf_config")
}

// GetClientLatestVersion 获取客户端最新版本
func (m *sysConfigModels) GetClientLatestVersion() string {
	return m.get("client_latest_version")
}

// GetClientLatestSourceURL 获取客户端最新下载地址
func (m *sysConfigModels) GetClientSourceURL() string {
	return m.get("client_source_url")
}
