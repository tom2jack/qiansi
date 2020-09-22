package mqttbroker

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zhi-miao/qiansi/common/config"

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

	// 设备上下线订阅
	sysClientOnlineSub = "$SYS/brokers/+/clients/+/+"
	// 注册通信主题
	regPub = "qiansi-client/reg/s/"
	regSub = "qiansi-client/reg/c/+"

	// 常规业务主题 设备用户，主题类型，发起方S/C
	defaultTopicFormat = "qiansi-client/chan/%s/%s/%s"

	// 初始化请求
	runInitTopic topic = "runInit"
	// 遥测
	telesignalTopic topic = "telesignal"
	// telegraf配置推送
	telegrafConfigTopic topic = "telegrafConfig"
	// 部署通道
	deployTopic topic = "deploy"
	// 监控指标订阅
	metricTopic topic = "metric"
	// 日志订阅
	logTopic topic = "log"
	// 客户端升级
	updateTopic topic = "update"
)

func sub() error {
	if token := mqttClient.Subscribe(sysClientOnlineSub, 0, clientOnlineStatusCallBack); token.Wait() && token.Error() != nil {
		return errors.New(sysClientOnlineSub)
	}
	if token := mqttClient.Subscribe(runInitTopic.Sub(), 0, runInitCallBack); token.Wait() && token.Error() != nil {
		return errors.New(runInitTopic.Sub())
	}
	if token := mqttClient.Subscribe(regSub, 0, registerCallBack); token.Wait() && token.Error() != nil {
		return errors.New(sysClientOnlineSub)
	}
	if token := mqttClient.Subscribe(deployTopic.Sub(), 0, deployCallBack); token.Wait() && token.Error() != nil {
		return errors.New(deployTopic.Sub())
	}
	if token := mqttClient.Subscribe(metricTopic.Sub(), 0, metricCallBack); token.Wait() && token.Error() != nil {
		return errors.New(metricTopic.Sub())
	}
	if token := mqttClient.Subscribe(logTopic.Sub(), 0, logCallBack); token.Wait() && token.Error() != nil {
		return errors.New(logTopic.Sub())
	}
	return nil
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
		logrus.Fatal("Can't connect mqtt broker!")
	}
	logrus.Info("mqtt service loading..")
	if err := sub(); err != nil {
		logrus.Fatal("mqtt broker Sub err", err.Error())
	}
}

func GetMqttClient() mqtt.Client {
	return mqttClient
}

type topic string

func (t topic) Pub(user string) string {
	return fmt.Sprintf(defaultTopicFormat, user, t, "S")
}

func (t topic) Sub() string {
	return fmt.Sprintf(defaultTopicFormat, "+", t, "C")
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
