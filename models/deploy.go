package models

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/utils"
	"time"
)

type Deploy struct {
	AfterCommand  string    `xorm:"not null comment('后置命令') VARCHAR(2000)"`
	BeforeCommand string    `xorm:"not null default '' comment('前置命令') VARCHAR(2000)"`
	Branch        string    `xorm:"not null default '' comment('git分支') VARCHAR(100)"`
	CreateTime    time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	DeployKeys    string    `xorm:"not null default '' comment('部署私钥') VARCHAR(3000)"`
	DeployType    int       `xorm:"not null default 0 comment('部署类型 0-本地 1-git 2-zip') TINYINT(1)"`
	Id            int       `xorm:"not null pk autoincr comment('应用唯一编号') index(IX_UID) INT(10)"`
	LocalPath     string    `xorm:"not null default '' comment('本地部署地址') VARCHAR(500)"`
	NowVersion    int       `xorm:"not null default 0 comment('当前版本') INT(10)"`
	OpenId        string    `xorm:"not null default '' unique CHAR(32)"`
	RemoteUrl     string    `xorm:"not null default '' comment('资源地址') VARCHAR(500)"`
	Title         string    `xorm:"not null default '' comment('应用名称') VARCHAR(120)"`
	Uid           int       `xorm:"not null comment('创建用户UID') index(IX_UID) INT(11)"`
	UpdateTime    time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

// List 获取应用列表
func (m *Deploy) List(offset int, limit int) ([]Deploy, int) {
	data := []Deploy{}
	rows := 0
	db := Mysql.Where("uid=?", m.Uid)
	if m.Title != "" {
		db = db.Where("title like ?", "%"+m.Title+"%")
	}
	db.Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

// BatchCheck 批量检测是否是当前用户应用
func (m *Deploy) BatchCheck(ids []int) bool {
	ids = utils.IdsFitter(ids)
	var count = 0
	db := Mysql.Model(m).Where("id in (?) and uid=?", ids, m.Uid).Count(&count)
	return db.Error == nil && count == len(ids)
}

// GetOpenId 获取openid
func (m *Deploy) GetOpenId() bool {
	db := Mysql.Select("open_id").Where("uid=?", m.Uid).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

func (m *Deploy) GetIdByOpenId() bool {
	db := Mysql.Select("id").Where("open_id=?", m.OpenId).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

// Count 统计当前用户部署应用数量
func (m *Deploy) Count() (num int, err error) {
	db := Mysql.Model(m).Where("uid=?", m.Uid).Count(&num)
	if db.Error != nil || db.RowsAffected == 0 {
		err = fmt.Errorf("查询失败")
	}
	return
}
