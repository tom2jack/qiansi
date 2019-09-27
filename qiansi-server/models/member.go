package models

import (
	"time"
)

type Member struct {
	CreateTime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	Id          int       `xorm:"not null pk autoincr INT(10)"`
	InviterUid  int       `xorm:"not null default 0 comment('邀请人UID') INT(10)"`
	MaxDeploy   int       `xorm:"not null default 15 comment('最大部署应用数量') INT(10)"`
	MaxSchedule int       `xorm:"not null default 15 comment('最大调度任务数量') INT(10)"`
	Password    string    `xorm:"not null comment('密码') VARCHAR(255)"`
	Phone       string    `xorm:"not null comment('手机号') unique CHAR(11)"`
	Status      int       `xorm:"not null default 1 comment('0-锁定 1-正常') TINYINT(1)"`
	UpdateTime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
}

// CheckDeploy 校验部署任务数量
func (m *Member) CheckDeploy() bool {
	var num int
	Mysql.Table("deploy").Where("uid=?", m.Id).Count(&num)
	Mysql.Select("max_deploy").First(m)
	if m.MaxDeploy == 0 || num == 0 {
		return false
	}
	return m.MaxDeploy > num
}

// CheckSchedule 校验调度任务数量
func (m *Member) CheckSchedule() bool {
	var num int
	Mysql.Table("schedule").Where("uid=?", m.Id).Count(&num)
	Mysql.Select("max_schedule").First(m)
	if m.MaxSchedule == 0 || num == 0 {
		return false
	}
	return m.MaxSchedule > num
}
