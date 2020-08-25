package models

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/zhi-miao/qiansi/common/utils"
)

type safeISMap struct {
	sync.RWMutex
	Map map[int]string
}

type serverModels struct {
	db *gorm.DB
}

var serverSecretCache *safeISMap

func init() {
	// 初始化serverSecretCache
	serverSecretCache = new(safeISMap)
	serverSecretCache.Map = make(map[int]string)
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
func (m *serverModels) Create(m *Server) error {
	return m.db.Create(m).Error
}

// GetApiSecret 获取服务器ApiSecret信息
func (m *Server) GetApiSecret(id int) (secret string) {
	serverSecretCache.RLock()
	secret = serverSecretCache.Map[id]
	serverSecretCache.RUnlock()
	if secret == "" {
		m.Id = id
		m.Get()
		serverSecretCache.Lock()
		serverSecretCache.Map[id] = m.ApiSecret
		serverSecretCache.Unlock()
		secret = m.ApiSecret
	}
	return
}

// Get 获取服务器信息
func (m *Server) Get() bool {
	db := Mysql.Model(m).Where("id=?", m.Id).First(m)
	return db.Error == nil && db.RowsAffected > 0
}

// List 获取应用列表
func (m *Server) List(offset int, limit int) ([]Server, int) {
	data := []Server{}
	rows := 0
	db := Mysql.Where("uid=?", m.Uid).Select("id, uid, server_name, server_status, server_rule_id, device_id, domain, create_time, update_time")
	if m.ServerName != "" {
		db = db.Where("server_name like ?", "%"+m.ServerName+"%")
	}
	db.Offset(offset).Limit(limit).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

func (m *Server) ListByUser() []Server {
	data := []Server{}
	Mysql.Where("uid = ?", m.Uid).Find(&data)
	return data
}

// BatchCheck 批量检测是否是当前用户服务器
func (m *Server) BatchCheck(ids []int) bool {
	ids = utils.IdsFitter(ids)
	var count = 0
	db := Mysql.Model(m).Where("id in (?) and uid=?", ids, m.Uid).Count(&count)
	return db.Error == nil && count == len(ids)
}

// Count 统计当前用户客戶端注冊數量
func (m *Server) Count() (num int, err error) {
	db := Mysql.Model(m).Where("uid=?", m.Uid).Count(&num)
	if db.Error != nil {
		err = fmt.Errorf("查询失败")
	}
	return
}

// UserServerIds 获取用户的服务器编号
func (m *Server) UserServerIds() (ids []int) {
	data := []Server{}
	Mysql.Model(m).Select("id").Where("uid = ?", m.Uid).Find(&data)
	for _, v := range data {
		ids = append(ids, v.Id)
	}
	return
}

//ExistsDevice 设备是否存在
func (m *serverModels) ExistsDevice(deviceID string) bool {
	dto := ModelBase{}
	err := m.db.Raw("select exists(select 1 from server where device_id=?) as has", deviceID).Scan(&dto).Error
	if err != nil {
		return true
	}
	return dto.Has
}
