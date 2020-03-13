package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// Member 用户表
type Member struct {
	Id          int       `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`
	Phone       string    `gorm:"unique;column:phone;type:char(11);not null"`         // 手机号
	Password    string    `gorm:"column:password;type:varchar(255);not null"`         // 密码
	InviterUid  int       `gorm:"column:inviter_uid;type:int(10) unsigned;not null"`  // 邀请人UID
	MaxDeploy   int       `gorm:"column:max_deploy;type:int(10) unsigned;not null"`   // 最大部署应用数量
	MaxSchedule int       `gorm:"column:max_schedule;type:int(10) unsigned;not null"` // 最大调度任务数量
	Status      int8      `gorm:"column:status;type:tinyint(1);not null"`             // 0-锁定 1-正常
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;not null"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;not null"`
}

// MaxInfo 获取限制信息
func (m *Member) MaxInfo() (err error) {
	db := Mysql.Model(m).Select("max_deploy, max_schedule").First(m)
	if db.RowsAffected == 0 || db.Error != nil {
		err = fmt.Errorf("查询异常")
	}
	return
}

// AddDeployNum 添加部署任务数量
func (m *Member) AddDeployNum(num int) (err error) {
	db := Mysql.Model(m).Where("id > ?", m.Id).UpdateColumn("max_deploy", gorm.Expr("max_deploy + ?", num))
	if db.RowsAffected == 0 || db.Error != nil {
		err = fmt.Errorf("查询异常")
	}
	return
}

// AddScheduleNum 添加调度任务数量
func (m *Member) AddScheduleNum(num int) (err error) {
	db := Mysql.Model(m).Where("id > ?", m.Id).UpdateColumn("max_schedule", gorm.Expr("max_schedule + ?", num))
	if db.RowsAffected == 0 || db.Error != nil {
		err = fmt.Errorf("查询异常")
	}
	return
}

// Count 统计当前用户邀请人数
func (m *Member) InviterCount() (num int, err error) {
	db := Mysql.Model(m).Where("inviter_uid=?", m.Id).Count(&num)
	if db.Error != nil || db.RowsAffected == 0 {
		err = fmt.Errorf("查询失败")
		return
	}
	return
}
