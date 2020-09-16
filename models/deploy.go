package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zhi-miao/gutils"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/zhi-miao/qiansi/common/req"
)

type deployModels struct {
	db *gorm.DB
}

// DeployTaskInfo 部署应用信息结构
type DeployTaskInfo struct {
	Deploy
	DeployGit    DeployGit    `json:"deploy_git"`
	DeployZip    DeployZip    `json:"deploy_zip"`
	DeployDocker DeployDocker `json:"deploy_docker"`
}

func GetDeployModels() *deployModels {
	return &deployModels{
		db: Mysql,
	}
}

func (m *deployModels) SetDB(db *gorm.DB) *deployModels {
	m.db = db
	return m
}

// UpdateServerDeployVersion 更新服务器部署版本
func (m *deployModels) UpdateServerDeployVersion(deployID, serverID, version int) error {
	return m.db.Model(&DeployServerRelation{}).
		Where("deploy_id=? and server_id=?", deployID, serverID).
		Update("deploy_version", version).Error
}

// UpdateVersion 更新当前版本
func (m *deployModels) UpdateVersion(UID, deployID int) error {
	tempDB := m.db.Model(&Deploy{}).Where("id=? and uid=?", deployID, UID).
		UpdateColumn("now_version", gorm.Expr("now_version+1"))
	if tempDB.Error != nil {
		return tempDB.Error
	}
	if tempDB.RowsAffected != 1 {
		return errors.New("部署任务不存在")
	}
	return nil
}

// GetDeployInfo 获取部署信息
func (m *deployModels) GetDeployInfo(deployID int) *DeployTaskInfo {
	result := &DeployTaskInfo{}
	if m.db.Model(&Deploy{}).Where("id=?", deployID).First(&result.Deploy).RowsAffected == 0 {
		return nil
	}
	// 部署类型 0-本地 1-git 2-zip 3-docker
	switch result.DeployType {
	case 0:
	case 1:
		m.db.Model(&DeployGit{}).Where("deploy_id=?", result.ID).First(&result.DeployGit)
	case 2:
		m.db.Model(&DeployZip{}).Where("deploy_id=?", result.ID).First(&result.DeployZip)
	case 3:
		m.db.Model(&DeployDocker{}).Where("deploy_id=?", result.ID).First(&result.DeployDocker)
	default:
	}
	return result
}

// GetDeployRelationServers 根据应用ID获取服务器信息
func (m *deployModels) GetDeployRelationServers(deployId int, serverIds ...int) []Server {
	sql := "select s.* from `server` s, deploy_server_relation r where s.id=r.server_id and r.deploy_id=?"
	data := make([]Server, 0)
	if len(serverIds) > 0 {
		m.db.Raw(sql+" and r.server_id in(?)", deployId, serverIds).Scan(&data)
	} else {
		m.db.Raw(sql, deployId).Scan(&data)
	}
	return data
}

// DeployList 部署应用列表
func (m *deployModels) DeployList(uid int, param *req.DeployListParam) (result []Deploy, totalRows int) {
	db := m.db.Model(&Deploy{}).Select("id, uid, title, deploy_type, now_version, open_id, create_time, update_time").Where("uid=?", uid)
	if param.Title != "" {
		db = db.Where("title like ?", "%"+param.Title+"%")
	}
	db.Offset(param.Offset()).Limit(param.PageSize).Order("id desc").Find(&result).Offset(-1).Limit(-1).Count(&totalRows)
	return
}

// GetDeployTaskInfo 根据服务器ID获取部署应用任务信息
func (m *deployModels) GetDeployTaskInfo(serverId int) (result []DeployTaskInfo) {
	result = []DeployTaskInfo{}
	m.db.Raw("SELECT d.* FROM `deploy` d LEFT JOIN `deploy_server_relation` r ON d.id=r.deploy_id WHERE r.server_id=? and d.now_version > r.deploy_version", serverId).Scan(&result)
	for i, info := range result {
		// 部署类型 0-本地 1-git 2-zip 3-docker
		switch info.DeployType {
		case 0:
		case 1:
			m.db.Model(&DeployGit{}).Where("deploy_id=?", info.ID).First(&result[i].DeployGit)
		case 2:
			m.db.Model(&DeployZip{}).Where("deploy_id=?", info.ID).First(&result[i].DeployZip)
		case 3:
			m.db.Model(&DeployDocker{}).Where("deploy_id=?", info.ID).First(&result[i].DeployDocker)
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
func (m *deployModels) GetDeployDetailInfo(uid, deployId int) (result DeployDetailInfo) {
	m.db.Model(&result.Deploy).Where("id=? and uid=?", deployId, uid).First(&result.Deploy)
	// 部署类型 0-本地 1-git 2-zip 3-docker
	switch result.DeployType {
	case 0:
	case 1:
		m.db.Model(&DeployGit{}).Where("deploy_id=?", result.ID).First(&result.DeployGit)
	case 2:
		m.db.Model(&DeployZip{}).Where("deploy_id=?", result.ID).First(&result.DeployZip)
	case 3:
		m.db.Model(&DeployDocker{}).Where("deploy_id=?", result.ID).First(&result.DeployDocker)
	default:
	}
	return
}

// CreateDeploy 创建部署应用
func (m *deployModels) CreateDeploy(uid int, param *req.DeploySetParam) (err error) {
	return m.db.Transaction(func(tx *gorm.DB) error {
		userInfo := &Member{Id: uid}
		info, err := GetMemberModels().UsereInfo(uid)
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
			// 检测服务器是否是自己的
			if !GetServerModels().SetDB(tx).BatchCheck(uid, serverIds) {
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
func (m *deployModels) UpdateDeploy(uid int, param *req.DeploySetParam) (err error) {
	return m.db.Transaction(func(tx *gorm.DB) error {
		info := Deploy{}
		if tx.Model(&info).Where("id=? and uid=?", param.ID, uid).First(&info).RowsAffected != 1 {
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
		switch info.DeployType {
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
			// 检测服务器是否是自己的
			if !GetServerModels().SetDB(tx).BatchCheck(uid, serverIds) {
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
func (m *deployModels) DelDeploy(uid, DeployID int) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		info := Deploy{ID: DeployID}
		db := tx.Model(&info).Where("uid=?", uid).First(&info)
		if db.Error != nil || db.RowsAffected == 0 {
			return fmt.Errorf("应用不存在")
		}
		tx.Delete(&info)
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

// DeployRelationServer 部署服务器信息
type DeployRelationServer struct {
	Server
	DeployServerRelation
}

// DeployServerList 可部署服务器列表
func (m *deployModels) DeployServerList(uid, deployId int) (result []DeployRelationServer, err error) {
	err = m.db.Raw("select a.*,b.deploy_id,b.deploy_version,a.id as server_id from `server` a LEFT JOIN (SELECT * from deploy_server_relation WHERE deploy_id=?) b ON a.id=b.server_id WHERE a.uid=?", deployId, uid).Scan(&result).Error
	return
}

// BatchCheck 批量检测是否是当前用户应用
func (m *deployModels) BatchCheck(ids []int, UID int) bool {
	ids = gutils.IdsUniqueFitter(ids)
	var count = 0
	db := m.db.Model(m).Where("id in (?) and uid=?", ids, UID).Count(&count)
	return db.Error == nil && count == len(ids)
}

// GetOpenID 获取openid
func (m *deployModels) GetOpenID(deployID, UID int) (string, error) {
	data := &Deploy{}
	db := m.db.Model(data).Select("open_id").Where("id=? and uid=?", deployID, UID).First(data)
	if db.Error != nil {
		return "", db.Error
	}
	return data.OpenID, nil
}

// GetIdByOpenId 根据openID换取部署应用ID
func (m *deployModels) GetIdByOpenId(openID string) (int, error) {
	d := &Deploy{}
	db := m.db.Model(d).Select("id").Where("open_id=?", openID).First(d)
	if db.Error != nil {
		return 0, db.Error
	}
	return d.ID, nil
}

// Count 统计当前用户部署应用数量
func (m *deployModels) Count(uid int) (int, error) {
	var num int
	err := m.db.Model(&Deploy{}).Where("uid=?", uid).Count(&num).Error
	return num, err
}

// Count 统计当前用户调度应用部署次数
func (m *deployModels) CountDo(uid int) (int, error) {
	r := TempModelStruct{}
	err := m.db.Model(&Deploy{}).Select("sum(now_version) as num").Where("uid=?", uid).Scan(&r).Error
	return r.Num, err
}
