package models

import (
	"crypto/sha256"
	"fmt"

	"github.com/jinzhu/gorm"
)

type mqttModels struct {
	db *gorm.DB
}

func GetMQTTModels() *mqttModels {
	return &mqttModels{
		db: Mysql,
	}
}

func (m *mqttModels) SetDB(db *gorm.DB) *mqttModels {
	m.db = db
	return m
}

// CreateClientUser 创建Mqtt客户端用户
func (m *mqttModels) CreateClientUser(username, password string) error {
	return m.db.Create(&MqttUser{
		Username:    username,
		Password:    fmt.Sprintf("%x", sha256.Sum256([]byte(password))),
		Salt:        "",
		IsSuperuser: false,
	}).Error
}

// CreateClientACL 创建客户端规则
func (m *mqttModels) CreateClientACL(username, clientID string) error {
	return m.db.Create(&MqttACL{
		Allow:    1,
		Username: username,
		Clientid: clientID,
		Access:   3,
		Topic:    fmt.Sprintf("qiansi-client/chan/%s/#", username),
	}).Error
}

// DeleteClientUser 删除客户端用户
func (m *mqttModels) DeleteClientUser(username string) error {
	return m.db.Delete(&MqttUser{}, "username=?", username).Error
}

// DeleteClientACL 删除客户端用户规则
func (m *mqttModels) DeleteClientACL(username string) error {
	return m.db.Delete(&MqttACL{}, "username=?", username).Error
}

//RestPassword 重置密码
func (m *mqttModels) RestPassword(username, newPwd string) error {
	return m.db.Model(&MqttUser{}).
		Where("username=?", username).
		Update("password", fmt.Sprintf("%x", sha256.Sum256([]byte(newPwd)))).Error
}
