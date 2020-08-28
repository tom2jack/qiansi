package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type memberModels struct {
	db *gorm.DB
}

func GetMemberModels() *memberModels {
	return &memberModels{
		db: Mysql,
	}
}

func (m *memberModels) SetDB(db *gorm.DB) *memberModels {
	m.db = db
	return m
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
	db := Mysql.Model(m).Where("inviter_uid=?", m.InviterUid).Count(&num)
	if db.Error != nil {
		err = fmt.Errorf("查询失败")
		return
	}
	return
}

//ExistsUID 判断uid是否存在
func (m *memberModels) ExistsUID(uid int) bool {
	dto := TempModelStruct{}
	err := m.db.Raw("select exists(select 1 from member where id=?) as has", uid).Scan(&dto).Error
	if err != nil {
		return true
	}
	return dto.Has
}
