package models

import (
	"errors"

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
func (m *memberModels) UsereInfo(uid int) (*Member, error) {
	userInfo := &Member{}
	db := m.db.Model(userInfo).Select("max_deploy, max_schedule").Where("id=?", uid).First(userInfo)
	if db.RowsAffected == 0 || db.Error != nil {
		return nil, errors.New("查询异常")
	}
	return userInfo, nil
}

// IncrMaxDeployNum 添加部署任务数量上限
func (m *memberModels) IncrMaxDeployNum(uid, num int) error {
	db := m.db.Model(m).Where("id=?", uid).UpdateColumn("max_deploy", gorm.Expr("max_deploy+?", num))
	if db.RowsAffected == 0 || db.Error != nil {
		return errors.New("提升失败")
	}
	return nil
}

// IncrMaxScheduleNum 添加调度任务数量上限
func (m *memberModels) IncrMaxScheduleNum(uid, num int) error {
	db := m.db.Model(m).Where("id=?", uid).UpdateColumn("max_schedule", gorm.Expr("max_schedule+?", num))
	if db.RowsAffected == 0 || db.Error != nil {
		return errors.New("提升失败")
	}
	return nil
}

// Count 统计当前用户邀请人数
func (m *memberModels) InviterCount(uid int) (int, error) {
	var num int
	err := m.db.Model(&Member{}).Where("inviter_uid=?", uid).Count(&num).Error
	return num, err
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
