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

func (m *Schedule) List(param *PageParam) (PageInfo, error) {
	data := []Schedule{}
	ZM_Mysql.Where("uid=?", m.Uid).Offset(param.Offset()).Limit(param.PageSize).Find(&data)
	// ZM_Mysql.Where("uid=?", m.Uid).
	// ZM_Mysql.Find()
	return PageInfo{
		Page: param.Page,
		PageSize: param.PageSize,
		TotalSize: 0,
		Rows: data,
	}, nil
}
