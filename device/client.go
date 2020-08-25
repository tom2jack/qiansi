package device

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/req"
)

// 客户端注册
func register(c mqtt.Client, message mqtt.Message) {
	defer message.Ack()
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
