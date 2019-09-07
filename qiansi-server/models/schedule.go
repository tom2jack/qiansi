package models

import (
	"time"
)

type Schedule struct {
	Command      string    `xorm:"not null comment('命令') TEXT"`
	CreateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	Crontab      string    `xorm:"not null comment('规则') VARCHAR(50)"`
	Id           int       `xorm:"not null pk autoincr INT(10)"`
	Remain       int       `xorm:"not null default -1 comment('剩余执行次数 -1无限 0-停止') INT(11)"`
	ScheduleType int       `xorm:"not null default 1 comment('调度类型 1-serverhttp 2-serverdeploy') TINYINT(1)"`
	Timeout      int       `xorm:"not null default 30 comment('超时时间s') INT(10)"`
	Title        string    `xorm:"not null comment('标题') VARCHAR(60)"`
	Uid          int       `xorm:"not null default 0 comment('用户ID') INT(10)"`
	UpdateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

// List 获取任务列表
func (m *Schedule) List(offset int, limit int) ([]Schedule, int) {
	data := []Schedule{}
	rows := 0
	Mysql.Where("uid=?", m.Uid).Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

// Create 创建任务
func (m *Schedule) Create() bool {
	if m.Uid <= 0 {
		return false
	}
	db := Mysql.Save(m)
	return db.Error == nil && db.RowsAffected > 0
}

// Del 删除计划任务
func (m *Schedule) Del() bool {
	db := Mysql.Delete(m, "id=? and uid=?", m.Id, m.Uid)
	return db.Error == nil && db.RowsAffected > 0
}
