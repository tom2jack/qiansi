package models

import (
	"fmt"
	"github.com/zhi-miao/qiansi/common/utils"
	"sync"
	"time"
)

type Server struct {
	ApiSecret    string    `xorm:"not null default '' comment('API密钥') VARCHAR(32)"`
	CreateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	DeviceId     string    `xorm:"not null comment('服务器唯一设备号') CHAR(36)"`
	Domain       string    `xorm:"not null comment('服务器地址(域名/ip)') VARCHAR(255)"`
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	ServerName   string    `xorm:"not null default '' comment('服务器备注名') VARCHAR(64)"`
	ServerRuleId int       `xorm:"not null default 0 comment('服务器规则id') INT(11)"`
	ServerStatus int       `xorm:"not null default 0 comment('服务器状态 -1-失效 0-待认领 1-已分配通信密钥 2-已绑定') TINYINT(1)"`
	Uid          int       `xorm:"not null comment('用户ID') index INT(10)"`
	UpdateTime   time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}

type safeISMap struct {
	sync.RWMutex
	Map map[int]string
}

var serverSecretCache *safeISMap

func init() {
	// 初始化serverSecretCache
	serverSecretCache = new(safeISMap)
	serverSecretCache.Map = make(map[int]string)
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
