package device

import (
	"encoding/json"
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lifei6671/gorand"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/models"
	"strings"
	"time"
)

// 客户端注册
func register(c mqtt.Client, message mqtt.Message) {
	defer message.Ack()
	type regResponse struct {
		ID        int    `json:"id"`
		Uid       int    `json:"uid"`
		DeviceId  string `json:"device_id"`
		ApiSecret string `json:"api_secret"`
		ErrMsg    string `json:"err_msg"`
	}
	type regParam struct {
		UID                   int    `json:"uid"`
		DeviceID              string `json:"device_id"`
		ResponseEncryptSecret string `json:"response_encrypt_secret"`
	}
	payload := message.Payload()
	param := regParam{}
	resp := regResponse{}
	err := func() error {
		err := json.Unmarshal(payload, &param)
		if err != nil {
			return errors.New("无法识别注册信息")
		}
		if !(param.UID > 0) {
			return errors.New("客户端唯一标识号非法")
		}
		member := &models.Member{}
		if !member.ExistsUID(param.UID) {
			return errors.New("用户不存在")
		}
		if len(param.DeviceID) != 36 {
			return errors.New("客户端唯一标识号非法")
		}
		server := models.Server{}
		if server.ExistsDevice(param.DeviceID) {
			return errors.New("设备已存在")
		}
		return nil
	}()
	if err == nil {
		server := &models.Server{
			Uid:          param.UID,
			ApiSecret:    string(gorand.KRand(16, gorand.KC_RAND_KIND_ALL)),
			DeviceId:     param.DeviceID,
			ServerStatus: 1,
		}
		err := server.Create()
		if err == nil {
			resp.ApiSecret = server.ApiSecret
			resp.ID = server.Id
		} else {
			resp.ErrMsg = "注册失败，请重新尝试"
		}
	} else {
		resp.ErrMsg = err.Error()
	}
	data, _ := json.Marshal(resp)
	raw := gutils.EncyptogAES(string(data), param.ResponseEncryptSecret)
	topicID := strings.TrimLeft(message.Topic(), topicRegSub[:len(topicRegSub)-1])
	c.Publish(topicRegPub+topicID, 0, false, raw).WaitTimeout(15 * time.Second)
}
