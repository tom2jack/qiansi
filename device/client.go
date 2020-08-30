package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/req"
	"github.com/zhi-miao/qiansi/service"
)

// registerCallBack 客户端注册
func registerCallBack(c mqtt.Client, message mqtt.Message) {
	payload := message.Payload()
	param := &req.RegServer{}
	err := json.Unmarshal(payload, param)
	if err != nil {
		logrus.Warn("报文无法识别", string(payload))
		return
	}
	data, err := json.Marshal(service.RegServer(param))
	if err != nil {
		logrus.Warn("客户端注册失败", err.Error())
		return
	}
	raw := gutils.EncyptogAES(string(data), param.ResponseEncryptSecret)
	topicID := strings.TrimLeft(message.Topic(), topicRegSub[:len(topicRegSub)-1])
	if !c.Publish(topicRegPub+topicID, 0, false, raw).WaitTimeout(waitTimeout) {
		logrus.Warn("mqtt服务器响应失败")
	}
}

// Deploy 启动部署
func Deploy(UID, deployID int, serverIds ...int) error {
	err := models.GetDeployModels().UpdateVersion(UID, deployID)
	if err != nil {
		return errors.New("版本更新失败")
	}
	deployModels := models.GetDeployModels()
	deployInfo := deployModels.GetDeployInfo(deployID)
	if deployInfo == nil {
		return errors.New("部署任务不存在")
	}
	payloadJson, err := json.Marshal(deployInfo)
	if err != nil {
		return errors.New("部署任务异常")
	}
	serversInfo := deployModels.GetDeployRelationServers(deployID, serverIds...)
	if len(serversInfo) == 0 {
		return errors.New("未挂载服务器")
	}
	for _, si := range serversInfo {
		mqttClient.Publish(fmt.Sprintf(topicDeployPub, si.MqttUser), 0, false, payloadJson)
	}
	return nil
}

// deployCallBack 部署回调
func deployCallBack(c mqtt.Client, message mqtt.Message) {
	info := GetAuthInfo(message.Topic())
	if info == nil {
		return
	}
	payload := &req.DeployCallBack{}
	err := UnmarshalPayload(message.Payload(), info.APISecret, payload)
	if err != nil {
		return
	}
	err = models.GetServerModels().UpdateDeployVersion(payload.DeployID, info.ServerID, payload.Version)
	if err != nil {
		return
	}
}
