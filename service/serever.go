package service

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lifei6671/gorand"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"
)

// UpdateServerOnlineStatus 更新设备上线状态
func UpdateServerOnlineStatus(deviceID string, isOnline bool) {

}

// RegServer 注册服务器
func RegServer(param *req.RegServer) (result *resp.RegServer) {
	result = &resp.RegServer{}
	if !(param.UID > 0) {
		result.ErrMsg = "客户端唯一标识号非法"
		return
	}
	if !models.GetMemberModels().ExistsUID(param.UID) {
		result.ErrMsg = "用户不存在"
		return
	}
	if len(param.DeviceID) != 36 {
		result.ErrMsg = "客户端唯一标识号非法"
		return
	}
	if models.GetServerModels().ExistsDevice(param.DeviceID) {
		result.ErrMsg = "设备已存在"
		return
	}
	err := models.Mysql.Transaction(func(tx *gorm.DB) error {
		server := &models.Server{
			UId:          param.UID,
			ServerStatus: 1,
			MqttUser:     "",
			APISecret:    string(gorand.KRand(16, gorand.KC_RAND_KIND_ALL)),
			DeviceID:     param.DeviceID,
			Domain:       "",
			CreateTime:   time.Time{},
			UpdateTime:   time.Time{},
		}
		sm := models.GetServerModels().SetDB(tx)
		err := sm.Create(server)
		if err != nil {
			return err
		}
		mm := models.GetMQTTModels().SetDB(tx)
		mqttPwd := string(gorand.KRand(16, gorand.KC_RAND_KIND_ALL))
		err = mm.CreateClientUser(server.MqttUser, mqttPwd)
		if err != nil {
			return err
		}
		err = mm.CreateClientACL(server.MqttUser, param.DeviceID)
		if err != nil {
			return err
		}
		result = &resp.RegServer{
			ID:               server.ID,
			ApiSecret:        server.APISecret,
			DeviceID:         param.DeviceID,
			MqttUserName:     server.MqttUser,
			MqttUserPassword: mqttPwd,
		}
		return nil
	})
	if err != nil {
		result.ErrMsg = "注册失败"
		return
	}
	return
}
