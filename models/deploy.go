package models

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/utils"
)

// List 获取应用列表
func (m *Deploy) List(offset int, limit int) ([]Deploy, int) {
	data := []Deploy{}
	rows := 0
	db := Mysql.Where("uid=?", m.UId)
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
	db := Mysql.Model(m).Where("id in (?) and uid=?", ids, m.UId).Count(&count)
	return db.Error == nil && count == len(ids)
}

// GetOpenId 获取openid
func (m *Deploy) GetOpenId() bool {
	db := Mysql.Select("open_id").Where("uid=?", m.UId).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

func (m *Deploy) GetIdByOpenId() bool {
	db := Mysql.Select("id").Where("open_id=?", m.OpenID).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

// Count 统计当前用户部署应用数量
func (m *Deploy) Count() (num int, err error) {
	db := Mysql.Model(m).Where("uid=?", m.UId).Count(&num)
	if db.Error != nil || db.RowsAffected == 0 {
		err = fmt.Errorf("查询失败")
	}
	return
}

// Count 统计当前用户调度应用部署次数
func (m *Deploy) CountDo() (num int, err error) {
	type Result struct {
		Total int
	}
	r := &Result{}
	db := Mysql.Model(m).Select("sum(now_version) as total").Where("uid=?", m.UId).Scan(r)
	if db.Error != nil || db.RowsAffected == 0 {
		err = fmt.Errorf("查询失败")
		return
	}
	num = r.Total
	return
}
