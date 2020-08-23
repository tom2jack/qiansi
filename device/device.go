package device

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common"
	"os"
	"time"
)

type service interface {
	Start()
}

var (
	mqttClient mqtt.Client
)

const (
	topicRegPub = "qiansi-client/reg/callback/"
	topicRegSub = "qiansi-client/reg/push/+"
)

func Start() {
	option := mqtt.NewClientOptions().
		AddBroker(common.Config.MQTT.Broker).
		SetAutoReconnect(true).
		SetUsername(common.Config.MQTT.Username).
		SetPassword(common.Config.MQTT.Password).
		SetClientID(common.Config.MQTT.ClientID)
	mqttClient = mqtt.NewClient(option)
	if !mqttClient.Connect().WaitTimeout(15 * time.Second) {
		logrus.Errorf("Can't connect mqtt broker!")
		os.Exit(1)
	}
	logrus.Info("mqtt service loading..")
	go sub()
}

func sub() {
	mqttClient.Subscribe(topicRegSub, 0, register)
}
