package models

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/req"
	"github.com/jinzhu/gorm"
)

// DeployList 部署应用列表
func DeployList(uid int, param *req.DeployListParam) (result []Deploy, totalRows int) {
	db := Mysql.Model(&Deploy{}).Select("id, uid, title, deploy_type, now_version, open_id, create_time, update_time").Where("uid=?", uid)
	if param.Title != "" {
		db = db.Where("title like ?", "%"+param.Title+"%")
	}
	db.Offset(param.Offset()).Limit(param.PageSize).Order("id desc").Find(&result).Offset(-1).Limit(-1).Count(&totalRows)
	return
}

// CreateDeploy 创建部署应用
func CreateDeploy(uid int, param req.DeploySetParam) (err error) {
	return
}

// UpdateDeploy 更新部署应用
func UpdateDeploy(uid int, param req.DeploySetParam) (err error) {
	return
}

// DelDeploy 删除部署应用
func DelDeploy(uid, DeployID int) error {
	return Mysql.Transaction(func(tx *gorm.DB) (err error) {
		info := Deploy{ID: DeployID}
		db := tx.Model(&info).Where("uid=?", uid).First(&info)
		if db.Error != nil || db.RowsAffected == 0 {
			err = fmt.Errorf("应用不存在")
			return
		}
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch info.DeployType {
		case 0:
		case 1:
			tx.Delete(&DeployGit{}, "deploy_id=?", DeployID)
		case 2:
			tx.Delete(&DeployZip{}, "deploy_id=?", DeployID)
		case 3:
			tx.Delete(&DeployDocker{}, "deploy_id=?", DeployID)
		default:
			return fmt.Errorf("无法识别的应用类型，请联系管理员删除")
		}
		// 删除服务器关联关系
		tx.Delete(&DeployServerRelation{}, "deploy_id=?", DeployID)
		return nil
	})
}

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
