package service

import (
	"errors"
	"fmt"
	"github.com/zhi-miao/qiansi/common/config"
	"github.com/zhi-miao/qiansi/common/sdk"
	"github.com/zhi-miao/qiansi/common/utils"
	"path"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lifei6671/gorand"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"
)

type serverService struct{}

func GetServerService() *serverService {
	return &serverService{}
}

// GetClientSource 根据系统架构获取客户端最新资源信息
func (s serverService) GetClientSource(os, arch string) (*resp.UpdateClient, error) {
	client, err := sdk.NewOSSClient()
	if err != nil {
		return nil, err
	}
	clientPath := config.GetConfig().Aliyun.OSS.QiansiClientScanPath
	file, err := client.ListFile(clientPath)
	if err != nil {
		return nil, err
	}
	for _, s := range file {
		_, f := path.Split(s)
		if f != "" {
			v, o, a, err := utils.ParseQiansiClientFileInfo(f)
			if err == nil && o == os && a == arch {
				ossConf := config.GetConfig().Aliyun.OSS
				return &resp.UpdateClient{
					Version:   v,
					SourceURL: fmt.Sprintf("https://%s/%s/%s", ossConf.Domain, ossConf.QiansiClientScanPath, f),
				}, nil
			}
		}
	}
	return nil, errors.New("not found the client source")
}

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
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
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
