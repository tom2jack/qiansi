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

// List 获取应用列表
func (m *DeployLog) List(start time.Time, end time.Time, offset int, limit int) ([]DeployLog, int) {
	data := []DeployLog{}
	rows := 0
	db := Mysql.Where("uid=? and (create_time between ? and ?)", m.Uid, start, end)
	if m.DeployId > 0 {
		db = db.Where("deploy_id=?", m.DeployId)
	}
	if m.DeployVersion > 0 {
		db = db.Where("deploy_version=?", m.DeployVersion)
	}
	if m.ServerId > 0 {
		db = db.Where("server_id=?", m.ServerId)
	}
	db.Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}
