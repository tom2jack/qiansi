package models

import (
	"github.com/jinzhu/gorm"
	"github.com/zhi-miao/qiansi/common/req"
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
func (m *scheduleModels) List(uid int, param *req.ScheduleListParam) ([]Schedule, int) {
	data := make([]Schedule, 0)
	rows := 0
	db := m.db.Where("uid=?", uid)
	if param.Title != "" {
		db = db.Where("title like ?", "%"+param.Title+"%")
	}
	db.Offset(param.Offset()).Limit(param.PageSize).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

func (m *scheduleModels) RunCallBack(param *Schedule) error {
	return m.db.Model(&Schedule{}).Where("id=?", param.ID).Updates(CommonMap{
		"prev_time": param.PrevTime,
		"next_time": param.NextTime,
		"remain":    param.Remain,
	}).Error
}

func (m *scheduleModels) Get(id, uid int) (Schedule, error) {
	data := Schedule{}
	err := m.db.Where("id=? and uid=?", id, uid).First(&data).Error
	return data, err
}

// Create 创建任务
func (m *scheduleModels) Create(param *Schedule) error {
	return m.db.Create(param).Error
}

// Del 删除计划任务
func (m *scheduleModels) Del(id, uid int) error {
	return m.db.Where("id=? and uid=?", id, uid).Delete(&Schedule{}).Error
}

func (m *scheduleModels) Save(param *Schedule) error {
	return m.db.Save(param).Error
}

// ExportList 数据输出
func (m *scheduleModels) ExportList(lastId int, scheduleType int) ([]Schedule, int) {
	data := []Schedule{}
	db := m.db
	if lastId > 0 {
		db = db.Where("id > ?", lastId)
	}
	if scheduleType > 0 {
		db = db.Where("schedule_type = ?", scheduleType)
	}
	db.Limit(1000).Order("id asc").Find(&data)
	return data, len(data)
}

// Count 统计当前用户调度应用数量
func (m *scheduleModels) Count(uid int) (int, error) {
	var num int
	err := m.db.Model(&Schedule{}).Where("uid=?", uid).Count(&num).Error
	return num, err
}
