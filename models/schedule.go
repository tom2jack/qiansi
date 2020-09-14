package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type scheduleModels struct {
	db *gorm.DB
}

func GetScheduleModels() *scheduleModels {
	return &scheduleModels{
		db: Mysql,
	}
}

func (m *scheduleModels) SetDB(db *gorm.DB) *scheduleModels {
	m.db = db
	return m
}

// List 获取任务列表
func (m *scheduleModels) List(offset int, limit int) ([]Schedule, int) {
	data := []Schedule{}
	rows := 0
	db := Mysql.Where("uid=?", m.Uid)
	if m.Title != "" {
		db = db.Where("title like ?", "%"+m.Title+"%")
	}
	db.Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

func (m *scheduleModels) Save() bool {
	db := Mysql.Save(m)
	return db.Error == nil && db.RowsAffected > 0
}

func (m *scheduleModels) RunCallBack() bool {
	db := Mysql.Model(m).Updates(CommonMap{"prev_time": m.PrevTime, "next_time": m.NextTime, "remain": m.Remain})
	return db.Error == nil && db.RowsAffected > 0
}

func (m *scheduleModels) Get() bool {
	db := Mysql.Where("id=? and uid=?", m.Id, m.Uid).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

// ExportList 数据输出
func (m *scheduleModels) ExportList(lastId int, scheduleType int) ([]Schedule, int) {
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
func (m *scheduleModels) Create() bool {
	if m.Uid <= 0 {
		return false
	}
	db := Mysql.Create(m)
	if db.Error != nil {
		logrus.Warnf("调度任务创建异常(rows: %d): %s", db.RowsAffected, db.Error)
	}
	return db.Error == nil && db.RowsAffected > 0
}

// Del 删除计划任务
func (m *scheduleModels) Del() bool {
	db := Mysql.Delete(m, "id=? and uid=?", m.Id, m.Uid)
	return db.Error == nil && db.RowsAffected > 0
}

// Count 统计当前用户调度应用数量
func (m *scheduleModels) Count() (num int, err error) {
	db := Mysql.Model(m).Where("uid=?", m.Uid).Count(&num)
	if db.Error != nil {
		err = fmt.Errorf("查询失败")
	}
	return
}
