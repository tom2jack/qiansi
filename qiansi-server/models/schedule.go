package models

import (
	"gitee.com/zhimiao/qiansi/common/logger"
	"time"
)

type Schedule struct {
	Command      string    `xorm:"not null comment('命令') TEXT"`
	CreateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	Crontab      string    `xorm:"not null comment('规则') VARCHAR(50)"`
	Id           int       `xorm:"not null pk autoincr INT(10)"`
	PrevTime     time.Time `gorm:"column:prev_time;type:datetime;default:null"` // 上次执行时间
	NextTime     time.Time `gorm:"column:next_time;type:datetime"`              // 下次执行时间
	Remain       int       `xorm:"not null default -1 comment('剩余执行次数 -1无限 0-停止') INT(11)"`
	ScheduleType int       `xorm:"not null default 1 comment('调度类型 1-serverhttp 2-clientShell') TINYINT(1)"`
	ServerId     int       `xorm:"not null default 0 comment('服务器ID') INT(10)"`
	Timeout      int       `xorm:"not null default 30 comment('超时时间s') INT(10)"`
	Title        string    `xorm:"not null comment('标题') VARCHAR(60)"`
	Uid          int       `xorm:"not null default 0 comment('用户ID') INT(10)"`
	UpdateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

// List 获取任务列表
func (m *Schedule) List(offset int, limit int) ([]Schedule, int) {
	data := []Schedule{}
	rows := 0
	db := Mysql.Where("uid=?", m.Uid)
	if m.Title != "" {
		db = db.Where("title like ?", "%"+m.Title+"%")
	}
	db.Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

func (m *Schedule) Save() bool {
	db := Mysql.Save(m)
	return db.Error == nil && db.RowsAffected > 0
}

func (m *Schedule) RunCallBack() bool {
	db := Mysql.Model(m).Updates(CommonMap{"prev_time": m.PrevTime, "next_time": m.NextTime, "remain": m.Remain})
	return db.Error == nil && db.RowsAffected > 0
}

func (m *Schedule) Get() bool {
	db := Mysql.Where("id=? and uid=?", m.Id, m.Uid).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

// ExportList 数据输出
func (m *Schedule) ExportList(lastId int, scheduleType int) ([]Schedule, int) {
	data := []Schedule{}
	db := Mysql
	if lastId > 0 {
		db = db.Where("id > ?", lastId)
	}
	if scheduleType > 0 {
		db = db.Where("schedule_type = ?", scheduleType)
	}
	db.Limit(1000).Order("id asc").Find(&data)
	return data, len(data)
}

// Create 创建任务
func (m *Schedule) Create() bool {
	if m.Uid <= 0 {
		return false
	}
	db := Mysql.Create(m)
	if db.Error != nil {
		logger.Warn("调度任务创建异常(rows: %d): %s", db.RowsAffected, db.Error)
	}
	return db.Error == nil && db.RowsAffected > 0
}

// Del 删除计划任务
func (m *Schedule) Del() bool {
	db := Mysql.Delete(m, "id=? and uid=?", m.Id, m.Uid)
	return db.Error == nil && db.RowsAffected > 0
}
