package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/errors"
	"github.com/zhi-miao/qiansi/common/req"
)

type serverModels struct {
	db    *gorm.DB
	redis *redis.Client
}

func GetServerModels() *serverModels {
	return &serverModels{
		db:    Mysql,
		redis: Redis,
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
	return errors.NewApiError(errors.Models, "创建失败")
}

// Get 获取服务器信息
func (m *serverModels) Get(serverID int) (*Server, error) {
	server := &Server{}
	err := m.db.Model(server).Where("id=?", serverID).First(server).Error
	return server, err
}

// GetByMqttUser 根据mqtt用户名获取服务器信息
func (m *serverModels) GetByMqttUser(mqttUser string) (*Server, error) {
	server := &Server{}
	err := m.db.Model(server).Where("mqtt_user=?", mqttUser).First(server).Error
	return server, err
}

// BatchCheck 批量检测是否是当前用户服务器
func (m *serverModels) BatchCheck(uid int, ids []int) bool {
	ids = gutils.IdsUniqueFitter(ids)
	var count int
	db := m.db.Model(&Server{}).Where("id in (?) and uid=?", ids, uid).Count(&count)
	return db.Error == nil && count == len(ids)
}

// UserServerIds 获取用户的服务器编号
func (m *serverModels) UserServerIds(UID int, tryCache bool) (ids []int) {
	if tryCache {
		val, err := m.redis.Get(context.Background(), fmt.Sprintf(UserServerIdsCacheKey, UID)).Result()
		if err == nil && val != "" {
			if json.Unmarshal([]byte(val), &ids) == nil {
				return
			}
		}
	}
	data := make([]Server, 0)
	m.db.Model(&Server{}).Select("id").Where("uid = ?", UID).Find(&data)
	for _, v := range data {
		ids = append(ids, v.ID)
	}
	if tryCache {
		jsonByte, err := json.Marshal(ids)
		if err == nil {
			m.redis.Set(context.Background(), fmt.Sprintf(UserServerIdsCacheKey, UID), string(jsonByte), 5*time.Minute)
		}
	}
	return
}

// GetServerIdByDeviceId 根据设备号换取服务器编号
func (m *serverModels) GetServerIdByDeviceId(deviceID string, tryCache bool) int {
	if tryCache {
		val, err := m.redis.HGet(context.Background(), ServerDeviceIDCacheKey, deviceID).Result()
		if err == nil {
			serverID, err := strconv.Atoi(val)
			if err == nil {
				return serverID
			}
		}
	}
	data := Server{}
	if m.db.Model(&Server{}).Select("id").Where("device_id=?", deviceID).First(&data).Error != nil {
		return 0
	}
	if tryCache {
		m.redis.HSet(context.Background(), ServerDeviceIDCacheKey, deviceID, data.ID)
	}
	return data.ID
}

// UpdateServerOnlineStat 更新设备在线状态
func (m *serverModels) UpdateServerOnlineStat(serverID int, isOnline bool) error {
	if isOnline {
		return m.redis.Set(context.Background(), fmt.Sprintf("%s:%d", ServerOnlineStatusCacheKey, serverID), time.Now().Unix(), 15*time.Second).Err()
	}
	return m.redis.Del(context.Background(), fmt.Sprintf("%s:%d", ServerOnlineStatusCacheKey, serverID)).Err()
}

// GetUserActiveServerNum 获取用户活跃服务器数
func (m *serverModels) GetUserActiveServerNum(UID int) int {
	serverIds := m.UserServerIds(UID, true)
	keys := make([]string, 0)
	for _, serverID := range serverIds {
		keys = append(keys, fmt.Sprintf("%s:%d", ServerOnlineStatusCacheKey, serverID))
	}
	statList, err := m.redis.MGet(context.Background(), keys...).Result()
	if err != nil {
		return 0
	}
	return len(statList)
}

// Count 统计当前用户客戶端注冊數量
func (m *serverModels) Count(UID int) int {
	var num int
	m.db.Model(&Server{}).Where("uid=?", UID).Count(&num)
	return num
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

// GetTelegrafConfig 获取Telegraf私有配置
func (m *serverModels) GetTelegrafConfig(serverID int) *Telegraf {
	data := &Telegraf{}
	m.db.Model(data).Where("server_id=?", serverID).First(data)
	return data
}

// List 获取服务器列表
func (m *serverModels) List(uid int, param *req.ServerListParam) ([]Server, int) {
	db := m.db.Model(&Server{}).
		Where("uid=?", uid).
		Select("id, uid, server_name, server_status, device_id, domain, create_time, update_time")
	if param.ServerName != "" {
		db = db.Where("server_name like ?", "%"+param.ServerName+"%")
	}
	var rows int
	data := make([]Server, 0)
	db.Offset(param.Offset()).Limit(param.PageSize).Order("id desc").Find(&data).Offset(-1).Limit(-1).Count(&rows)
	return data, rows
}

// ListByUID 根据用户ID获取服务器列表
func (m *serverModels) ListByUID(uid int) ([]Server, error) {
	data := []Server{}
	err := m.db.Model(&Server{}).Where("uid = ?", uid).Find(&data).Error
	return data, err
}

// Delete 根据用户ID删除服务器
func (m *serverModels) Delete(id, uid int) error {
	servInfo, err := m.Get(id)
	if err != nil {
		return err
	}
	if servInfo.UId != uid {
		return errors.NewApiError(errors.Models, "非本人服务器无权删除")
	}
	db := m.db.Delete(Server{}, "id=? and uid=?", id, uid)
	if db.Error != nil || db.RowsAffected != 1 {
		return errors.NewApiError(errors.Models, "服务器删除失败")
	}
	// 删除关联
	if m.db.Delete(DeployServerRelation{}, "server_id=?", id).Error != nil {
		return errors.NewApiError(errors.Models, "服务器部署关联删除失败")
	}
	// 删除mqtt用户
	mqttModel := GetMQTTModels().SetDB(m.db)
	err = mqttModel.DeleteClientUser(servInfo.MqttUser)
	if err != nil {
		return err
	}
	// 删除mqtt用户权限
	err = mqttModel.DeleteClientACL(servInfo.MqttUser)
	if err != nil {
		return err
	}
	return nil
}
