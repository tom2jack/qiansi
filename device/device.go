package device

import (
	"encoding/json"
	"errors"
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/models"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common"
)

var (
	mqttClient mqtt.Client
)

const (
	// 注册通信主题
	topicRegPub = "qiansi-client/reg/s/"
	topicRegSub = "qiansi-client/reg/c/+"
	// 部署通道
	topicDeployPub = "qiansi-client/chan/%s/deploy/S"
	topicDeploySub = "qiansi-client/chan/+/deploy/C"

	waitTimeout = 15 * time.Second
)

func sub() {
	mqttClient.Subscribe(topicRegSub, 0, registerCallBack)
	mqttClient.Subscribe(topicDeploySub, 0, deployCallBack)
}

func Start() {
	option := mqtt.NewClientOptions().
		AddBroker(common.Config.MQTT.Broker).
		SetAutoReconnect(true).
		SetUsername(common.Config.MQTT.Username).
		SetPassword(common.Config.MQTT.Password).
		SetClientID(common.Config.MQTT.ClientID)
	mqttClient = mqtt.NewClient(option)
	if !mqttClient.Connect().WaitTimeout(waitTimeout) {
		logrus.Errorf("Can't connect mqtt broker!")
		os.Exit(1)
	}
	logrus.Info("mqtt service loading..")
	go sub()
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

func GetAuthInfo(topic string) *authInfo {
	topicSplit := strings.Split(topic, "/")
	if len(topicRegSub) < 3 {
		return nil
	}
	auth := &authInfo{}
	auth.UserName = topicSplit[2]
	auth.ServerID, _ = strconv.Atoi(strings.TrimLeft(auth.UserName, "Q_"))
	once, err := models.GetServerModels().GetOnce(auth.ServerID)
	if err != nil {
		return nil
	}
	auth.APISecret = once.APISecret
	auth.DeviceID = once.DeviceID
	auth.UID = once.UId
	return auth
}

// UnmarshalPayload 解码载荷
func UnmarshalPayload(payload []byte, secret string, res interface{}) error {
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

// MarshalPayload 加密载荷
func MarshalPayload(payload interface{}, secret string) ([]byte, error) {
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
