package service

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lifei6671/gorand"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/req"
	"github.com/zhi-miao/qiansi/resp"
)

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
			Uid:          param.UID,
			ApiSecret:    string(gorand.KRand(16, gorand.KC_RAND_KIND_ALL)),
			DeviceId:     param.DeviceID,
			ServerStatus: 1,
		}
		sm := models.GetServerModels().SetDB(tx)
		err := sm.Create(server)
		if err != nil {
			return err
		}
		mm := models.GetMQTTModels().SetDB(tx)
		username := "Q_"+ server.Id
		mm.CreateClientUser()
	})
	if err != nil {
		result.ErrMsg = "注册失败"
		return
	}
	err := server.Create()
	if err == nil {
		resp.ApiSecret = server.ApiSecret
		resp.ID = server.Id
	} else {

	}
}
