package models

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/req"
	"gitee.com/zhimiao/qiansi/service"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"strings"
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
func CreateDeploy(uid int, param *req.DeploySetParam) (err error) {
	return Mysql.Transaction(func(tx *gorm.DB) error {
		info, err := service.GetUserModuleMaxInfo(uid)
		if err != nil {
			return fmt.Errorf("用户限额检测失败")
		}
		if info.MaxDeploy <= info.DeployNum {
			return fmt.Errorf("您的部署任务创建数量已达上限")
		}
		deploy := Deploy{
			UId:           uid,
			Title:         param.Title,
			DeployType:    param.DeployType,
			WorkDir:       param.WorkDir,
			BeforeCommand: param.BeforeCommand,
			AfterCommand:  param.AfterCommand,
			NowVersion:    0,
			OpenID:        strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		}
		err = tx.Save(&deploy).Error
		if err != nil {
			return fmt.Errorf("应用创建失败")
		}
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch deploy.DeployType {
		case 0:
		case 1:
			err = tx.Save(&DeployGit{
				DeployID:   deploy.ID,
				RemoteURL:  param.DeployGit.RemoteURL,
				DeployPath: param.DeployGit.DeployPath,
				Branch:     param.DeployGit.Branch,
				DeployKeys: param.DeployGit.DeployKeys,
				UserName:   param.DeployGit.UserName,
				Password:   param.DeployGit.Password,
			}).Error
		case 2:
			err = tx.Save(&DeployZip{
				DeployID:   deploy.ID,
				RemoteURL:  param.DeployZip.RemoteURL,
				DeployPath: param.DeployZip.DeployPath,
				Password:   param.DeployZip.Password,
			}).Error
		case 3:
			err = tx.Save(&DeployDocker{
				DeployID:         deploy.ID,
				DockerImage:      param.DeployDocker.DockerImage,
				UserName:         param.DeployDocker.UserName,
				Password:         param.DeployDocker.Password,
				IsRuning:         param.DeployDocker.IsRuning,
				ContainerName:    param.DeployDocker.ContainerName,
				ContainerVolumes: param.DeployDocker.ContainerVolumes,
				ContainerPorts:   param.DeployDocker.ContainerPorts,
				ContainerEnv:     param.DeployDocker.ContainerEnv,
			}).Error
		default:
			return fmt.Errorf("无法识别的应用类型")
		}
		if err != nil {
			return fmt.Errorf("应用创建失败")
		}
		serverIds := make([]int, len(param.ServerRelation))
		for _, s := range param.ServerRelation {
			serverIds = append(serverIds, s.ServerId)
		}
		server := &Server{Uid: uid}
		if !server.BatchCheck(serverIds) {
			return fmt.Errorf("含有非法服务器绑定")
		}
		relation := &DeployServerRelation{}
		for _, e := range param.ServerRelation {
			relation.DeployID = deploy.ID
			relation.ServerID = e.ServerId
			if tx.Save(relation).Error != nil {
				return fmt.Errorf("服务器关联失败")
			}
		}
		return nil
	})
}

// UpdateDeploy 更新部署应用
func UpdateDeploy(uid int, param *req.DeploySetParam) (err error) {
	// TODO:
	return Mysql.Transaction(func(tx *gorm.DB) error {
		info, err := service.GetUserModuleMaxInfo(uid)
		if err != nil {
			return fmt.Errorf("用户限额检测失败")
		}
		if info.MaxDeploy <= info.DeployNum {
			return fmt.Errorf("您的部署任务创建数量已达上限")
		}
		deploy := Deploy{
			ID:            param.ID,
			UId:           uid,
			Title:         param.Title,
			DeployType:    param.DeployType,
			WorkDir:       param.WorkDir,
			BeforeCommand: param.BeforeCommand,
			AfterCommand:  param.AfterCommand,
			NowVersion:    0,
			OpenID:        strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
		}
		err = tx.Model(&deploy).Update(&deploy).Error
		if err != nil {
			return fmt.Errorf("应用创建失败")
		}
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch deploy.DeployType {
		case 0:
		case 1:
			err = tx.Model(&DeployGit{}).Update(&DeployGit{
				DeployID:   deploy.ID,
				RemoteURL:  param.DeployGit.RemoteURL,
				DeployPath: param.DeployGit.DeployPath,
				Branch:     param.DeployGit.Branch,
				DeployKeys: param.DeployGit.DeployKeys,
				UserName:   param.DeployGit.UserName,
				Password:   param.DeployGit.Password,
			}).Error
		case 2:
			err = tx.Model(&DeployZip{}).Update(&DeployZip{
				DeployID:   deploy.ID,
				RemoteURL:  param.DeployZip.RemoteURL,
				DeployPath: param.DeployZip.DeployPath,
				Password:   param.DeployZip.Password,
			}).Error
		case 3:
			err = tx.Model(&DeployDocker{}).Update(&DeployDocker{
				DeployID:         deploy.ID,
				DockerImage:      param.DeployDocker.DockerImage,
				UserName:         param.DeployDocker.UserName,
				Password:         param.DeployDocker.Password,
				IsRuning:         param.DeployDocker.IsRuning,
				ContainerName:    param.DeployDocker.ContainerName,
				ContainerVolumes: param.DeployDocker.ContainerVolumes,
				ContainerPorts:   param.DeployDocker.ContainerPorts,
				ContainerEnv:     param.DeployDocker.ContainerEnv,
			}).Error
		default:
			return fmt.Errorf("无法识别的应用类型")
		}
		if err != nil {
			return fmt.Errorf("应用创建失败")
		}
		serverIds := make([]int, len(param.ServerRelation))
		for _, s := range param.ServerRelation {
			serverIds = append(serverIds, s.ServerId)
		}
		server := &Server{Uid: uid}
		if !server.BatchCheck(serverIds) {
			return fmt.Errorf("含有非法服务器绑定")
		}
		relation := &DeployServerRelation{}
		for _, e := range param.ServerRelation {
			relation.DeployID = deploy.ID
			relation.ServerID = e.ServerId
			if tx.Save(relation).Error != nil {
				return fmt.Errorf("服务器关联失败")
			}
		}
		return nil
	})
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
