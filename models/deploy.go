package models

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/utils"
	"gitee.com/zhimiao/qiansi/req"
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

// DeployTaskInfo 部署应用信息结构
type DeployTaskInfo struct {
	Deploy
	DeployGit    DeployGit    `json:"deploy_git"`
	DeployZip    DeployZip    `json:"deploy_zip"`
	DeployDocker DeployDocker `json:"deploy_docker"`
}

// GetDeployTaskInfo 根据服务器ID获取部署应用任务信息
func GetDeployTaskInfo(serverId int) (result []DeployTaskInfo) {
	result = []DeployTaskInfo{}
	Mysql.Raw("SELECT d.* FROM `deploy` d LEFT JOIN `deploy_server_relation` r ON d.id=r.deploy_id WHERE r.server_id=? and d.now_version > r.deploy_version", serverId).Scan(&result)
	for i, info := range result {
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch info.DeployType {
		case 0:
		case 1:
			Mysql.Model(&DeployGit{}).Where("deploy_id=?", info.ID).First(&result[i].DeployGit)
		case 2:
			Mysql.Model(&DeployZip{}).Where("deploy_id=?", info.ID).First(&result[i].DeployZip)
		case 3:
			Mysql.Model(&DeployDocker{}).Where("deploy_id=?", info.ID).First(&result[i].DeployDocker)
		default:
		}
	}
	return
}

// DeployDetailInfo 部署应用信息结构
type DeployDetailInfo struct {
	Deploy
	DeployGit    DeployGit    `json:"deploy_git"`
	DeployZip    DeployZip    `json:"deploy_zip"`
	DeployDocker DeployDocker `json:"deploy_docker"`
}

// GetDeployDetailInfo 根据服应用ID获取部署应用任务信息
func GetDeployDetailInfo(uid, deployId int) (result DeployDetailInfo) {
	Mysql.Model(&result.Deploy).Where("id=? and uid=?", deployId, uid).First(&result.Deploy)
	// 部署类型 0-本地 1-git 2-zip 3-docker
	switch result.DeployType {
	case 0:
	case 1:
		Mysql.Model(&DeployGit{}).Where("deploy_id=?", result.ID).First(&result.DeployGit)
	case 2:
		Mysql.Model(&DeployZip{}).Where("deploy_id=?", result.ID).First(&result.DeployZip)
	case 3:
		Mysql.Model(&DeployDocker{}).Where("deploy_id=?", result.ID).First(&result.DeployDocker)
	default:
	}
	return
}

// CreateDeploy 创建部署应用
func CreateDeploy(uid int, param *req.DeploySetParam) (err error) {
	return Mysql.Transaction(func(tx *gorm.DB) error {
		userInfo := &Member{Id: uid}
		err := userInfo.MaxInfo()
		if err != nil {
			return fmt.Errorf("用户限额检测失败")
		}
		var deployNum int
		err = tx.Model(&Deploy{}).Where("uid=?", uid).Count(&deployNum).Error
		if err != nil {
			return fmt.Errorf("部署应用数量查询失败")
		}
		if userInfo.MaxDeploy <= deployNum {
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
		err = tx.Model(&deploy).Create(&deploy).Error
		if err != nil {
			return fmt.Errorf("应用创建失败")
		}
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch deploy.DeployType {
		case 0:
		case 1:
			err = tx.Model(&DeployGit{}).Create(&DeployGit{
				DeployID:   deploy.ID,
				RemoteURL:  param.DeployGit.RemoteURL,
				DeployPath: param.DeployGit.DeployPath,
				Branch:     param.DeployGit.Branch,
				DeployKeys: param.DeployGit.DeployKeys,
				UserName:   param.DeployGit.UserName,
				Password:   param.DeployGit.Password,
			}).Error
		case 2:
			err = tx.Model(&DeployZip{}).Create(&DeployZip{
				DeployID:   deploy.ID,
				RemoteURL:  param.DeployZip.RemoteURL,
				DeployPath: param.DeployZip.DeployPath,
				Password:   param.DeployZip.Password,
			}).Error
		case 3:
			err = tx.Model(&DeployDocker{}).Create(&DeployDocker{
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
		if len(param.ServerRelation) > 0 {
			serverIds := make([]int, len(param.ServerRelation))
			for i, s := range param.ServerRelation {
				serverIds[i] = s.ServerId
			}
			server := &Server{Uid: uid}
			if !server.BatchCheck(serverIds) {
				return fmt.Errorf("含有非法服务器绑定")
			}
			relation := &DeployServerRelation{}
			for _, e := range param.ServerRelation {
				if e.ServerId == 0 {
					continue
				}
				relation.DeployID = deploy.ID
				relation.ServerID = e.ServerId
				if tx.Model(&DeployServerRelation{}).Create(relation).Error != nil {
					return fmt.Errorf("服务器关联失败")
				}
			}
		}
		return nil
	})
}

// UpdateDeploy 更新部署应用
func UpdateDeploy(uid int, param *req.DeploySetParam) (err error) {
	return Mysql.Transaction(func(tx *gorm.DB) error {
		scan := ModelBase{}
		tx.Raw("SELECT EXISTS(SELECT 1 FROM deploy WHERE id=? and uid=?) as has", param.ID, uid).Scan(&scan)
		if !scan.Has {
			return fmt.Errorf("应用不存在")
		}
		deploy := Deploy{
			ID:            param.ID,
			Title:         param.Title,
			WorkDir:       param.WorkDir,
			BeforeCommand: param.BeforeCommand,
			AfterCommand:  param.AfterCommand,
		}
		err = tx.Model(&deploy).Update(&deploy).Error
		if err != nil {
			return fmt.Errorf("应用更新失败")
		}
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch deploy.DeployType {
		case 0:
		case 1:
			err = tx.Model(&DeployGit{}).Update(&DeployGit{
				DeployID:   param.ID,
				RemoteURL:  param.DeployGit.RemoteURL,
				DeployPath: param.DeployGit.DeployPath,
				Branch:     param.DeployGit.Branch,
				DeployKeys: param.DeployGit.DeployKeys,
				UserName:   param.DeployGit.UserName,
				Password:   param.DeployGit.Password,
			}).Error
		case 2:
			err = tx.Model(&DeployZip{}).Update(&DeployZip{
				DeployID:   param.ID,
				RemoteURL:  param.DeployZip.RemoteURL,
				DeployPath: param.DeployZip.DeployPath,
				Password:   param.DeployZip.Password,
			}).Error
		case 3:
			err = tx.Model(&DeployDocker{}).Update(&DeployDocker{
				DeployID:         param.ID,
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
			return fmt.Errorf("应用更新失败")
		}
		if len(param.ServerRelation) > 0 {
			serverIds := make([]int, len(param.ServerRelation))
			for i, s := range param.ServerRelation {
				serverIds[i] = s.ServerId
			}
			server := &Server{Uid: uid}
			if !server.BatchCheck(serverIds) {
				return fmt.Errorf("含有非法服务器绑定")
			}
			relation := &DeployServerRelation{}
			for _, e := range param.ServerRelation {
				if e.ServerId == 0 {
					continue
				}
				relation.DeployID = deploy.ID
				relation.ServerID = e.ServerId
				if e.Relation {
					if tx.Save(relation).Error != nil {
						return fmt.Errorf("服务器关联失败")
					}
				} else if tx.Delete(relation, "server_id=? and deploy_id=?", relation.ServerID, relation.DeployID).Error != nil {
					return fmt.Errorf("服务器取消关联失败")
				}
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

// DeployRelationServerIds 根据应用ID
func DeployRelationServerIds(deployId int) []int {
	data := []DeployServerRelation{}
	Mysql.Where("deploy_id=?", deployId).Find(&data)
	result := make([]int, len(data))
	for _, datum := range data {
		result = append(result, datum.ServerID)
	}
	return result
}

// DeployRelationServer 部署服务器信息
type DeployRelationServer struct {
	Server
	DeployServerRelation
}

// DeployServerList 可部署服务器列表
func DeployServerList(uid, deployId int) (result []DeployRelationServer, err error) {
	err = Mysql.Raw("select a.*,b.deploy_id,b.deploy_version from `server` a LEFT JOIN (SELECT * from deploy_server_relation WHERE deploy_id=?) b ON a.id=b.server_id WHERE a.uid=?", deployId, uid).Scan(&result).Error
	return
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
