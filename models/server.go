package models

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/zhi-miao/gutils"
)

type safeISMap struct {
	sync.RWMutex
	Map map[int]string
}

type serverModels struct {
	db *gorm.DB
}

func GetServerModels() *serverModels {
	return &serverModels{
		db: Mysql,
	}
}

func (m *serverModels) SetDB(db *gorm.DB) *serverModels {
	m.db = db
	return m
}

// Create 创建客户端
func (m *serverModels) Create(ser *Server) error {
	err := m.db.Create(ser).Error
	if err != nil {
		return err
	}
	if ser.ID > 0 {
		ser.MqttUser = fmt.Sprintf("Q_%d", ser.ID)
		err := m.db.Model(ser).Where("id=?", ser.ID).Update("mqtt_user", ser.MqttUser).Error
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("创建失败")
}

// Get 获取服务器信息
func (m *serverModels) GetOnce(serverID int) (*Server, error) {
	server := &Server{}
	err := m.db.Model(server).Where("id=?", serverID).First(server).Error
	return server, err
}

// BatchCheck 批量检测是否是当前用户服务器
func (m *serverModels) BatchCheck(uid int, ids []int) bool {
	ids = gutils.IdsUniqueFitter(ids)
	var count int
	db := m.db.Model(&Server{}).Where("id in (?) and uid=?", ids, uid).Count(&count)
	return db.Error == nil && count == len(ids)
}

func (m *serverModels) UpdateDeployVersion(deployID, serverID, version int) error {
	return m.db.Model(&DeployServerRelation{}).
		Where("deploy_id=? and server_id=?", deployID, serverID).
		Update("deploy_version", version).Error
}

// UserServerIds 获取用户的服务器编号
func (m *serverModels) UserServerIds(UID int) (ids []int) {
	data := make([]Server, 0)
	Mysql.Model(&Server{}).Select("id").Where("uid = ?", UID).Find(&data)
	for _, v := range data {
		ids = append(ids, v.ID)
	}
	return
}

// Count 统计当前用户客戶端注冊數量
func (m *serverModels) Count(UID int) int {
	var num int
	m.db.Model(&Server{}).Where("uid=?", UID).Count(&num)
	return num
}

// List 获取应用列表
func (m *Server) List(offset int, limit int) ([]Server, int) {
	data := []Server{}
	rows := 0
	db := Mysql.Where("uid=?", m.UId).Select("id, uid, server_name, server_status, server_rule_id, device_id, domain, create_time, update_time")
	if m.ServerName != "" {
		db = db.Where("server_name like ?", "%"+m.ServerName+"%")
	}
	db.Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

func (m *Server) ListByUser() []Server {
	data := []Server{}
	Mysql.Where("uid = ?", m.UId).Find(&data)
	return data
}

//ExistsDevice 设备是否存在
func (m *serverModels) ExistsDevice(deviceID string) bool {
	dto := TempModelStruct{}
	err := m.db.Raw("select exists(select 1 from server where device_id=?) as has", deviceID).Scan(&dto).Error
	if err != nil {
		return true
	}
	return dto.Has
}
