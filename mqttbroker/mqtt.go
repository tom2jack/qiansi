package mqttbroker

import (
	"encoding/json"
	"errors"
	"github.com/zhi-miao/qiansi/common/config"
	"os"
	"strings"
	"time"

	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var (
	mqttClient mqtt.Client
)

const (
	waitTimeout = 15 * time.Second
)

const (
	// 设备上下线订阅
	sysClientOnlineSub = "$SYS/brokers/+/clients/+/+"
	// 注册通信主题
	regPub = "qiansi-client/reg/s/"
	regSub = "qiansi-client/reg/c/+"
	// 初始化请求
	runInitSub = "qiansi-client/chan/+/runInit/C"
	// telegraf配置推送
	telegrafConfigPub = "qiansi-client/chan/%s/telegrafConfig/S"
	// 部署通道
	deployPub = "qiansi-client/chan/%s/deploy/S"
	deploySub = "qiansi-client/chan/+/deploy/C"
	// 监控指标订阅
	metricSub = "qiansi-client/chan/+/metric/C"
	// 日志订阅
	logSub = "qiansi-client/chan/+/log/C"
)

func sub() {
	if token := mqttClient.Subscribe(sysClientOnlineSub, 0, clientOnlineStatusCallBack); token.Wait() && token.Error() != nil {
		logrus.Warn("订阅失败->", sysClientOnlineSub)
		return
	}
	if token := mqttClient.Subscribe(runInitSub, 0, runInitCallBack); token.Wait() && token.Error() != nil {
		logrus.Warn("订阅失败->", runInitSub)
		return
	}
	if token := mqttClient.Subscribe(regSub, 0, registerCallBack); token.Wait() && token.Error() != nil {
		logrus.Warn("订阅失败->", regSub)
		return
	}
	if token := mqttClient.Subscribe(deploySub, 0, deployCallBack); token.Wait() && token.Error() != nil {
		logrus.Warn("订阅失败->", deploySub)
		return
	}
	if token := mqttClient.Subscribe(metricSub, 0, metricCallBack); token.Wait() && token.Error() != nil {
		logrus.Warn("订阅失败->", metricSub)
		return
	}
	if token := mqttClient.Subscribe(logSub, 0, logCallBack); token.Wait() && token.Error() != nil {
		logrus.Warn("订阅失败->", logSub)
		return
	}
}

func Start() {
	option := mqtt.NewClientOptions().
		AddBroker(config.GetConfig().MQTT.Broker).
		SetAutoReconnect(true).
		SetUsername(config.GetConfig().MQTT.Username).
		SetPassword(config.GetConfig().MQTT.Password).
		SetClientID(config.GetConfig().MQTT.ClientID)
	mqttClient = mqtt.NewClient(option)
	if !mqttClient.Connect().WaitTimeout(waitTimeout) {
		logrus.Errorf("Can't connect mqtt broker!")
		os.Exit(1)
	}
	logrus.Info("mqtt service loading..")
	sub()
}

func GetMqttClient() mqtt.Client {
	return mqttClient
}

type authInfo struct {
	UID       int
	ServerID  int
	DeviceID  string
	UserName  string
	APISecret string
}

// getAuthInfo 获取消息用户信息
func getAuthInfo(topic string) *authInfo {
	topicSplit := strings.Split(topic, "/")
	if len(topicSplit) < 3 {
		return nil
	}
	sInfo, err := models.GetServerModels().GetByMqttUser(topicSplit[2])
	if err != nil {
		return nil
	}
	return &authInfo{
		UID:       sInfo.UId,
		ServerID:  sInfo.ID,
		DeviceID:  sInfo.DeviceID,
		UserName:  sInfo.MqttUser,
		APISecret: sInfo.APISecret,
	}
}

// dePayload 解码载荷
func dePayload(payload []byte, secret string, res interface{}) error {
	raw := gutils.DecrptogAES(string(payload), secret)
	if raw == "" {
		return errors.New("解码失败")
	}
	err := json.Unmarshal([]byte(raw), res)
	if err != nil {
		return err
	}
	return nil
}

// enPayload 加密载荷
func enPayload(payload interface{}, secret string) ([]byte, error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	aes := gutils.EncyptogAES(string(payloadJson), secret)
	if aes == "" {
		return nil, errors.New("加密失败")
	}
	return []byte(aes), nil
}
