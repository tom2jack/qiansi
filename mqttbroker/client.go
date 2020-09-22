package mqttbroker

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/zhi-miao/qiansi/common/resp"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/service"
)

// clientOnlineStatusCallBack 客户端上下线状态变更回调
func clientOnlineStatusCallBack(c mqtt.Client, message mqtt.Message) {
	topicArr := strings.Split(message.Topic(), "/")
	arrLen := len(topicArr)
	stat := topicArr[arrLen-1]
	clientID := topicArr[arrLen-2]
	serverModel := models.GetServerModels()
	serverID := serverModel.GetServerIdByDeviceId(clientID, true)
	if serverID == 0 {
		return
	}
	//  connected 上线, disconnected 下线
	err := serverModel.UpdateServerOnlineStat(serverID, stat == "connected")
	if err != nil {
		logrus.Warn("设备上下线更新失败", err.Error())
	}
}

// registerCallBack 客户端注册
func registerCallBack(c mqtt.Client, message mqtt.Message) {
	payload := message.Payload()
	param := &req.RegServer{}
	err := json.Unmarshal(payload, param)
	if err != nil {
		logrus.Warn("报文无法识别", string(payload))
		return
	}
	resp := service.GetServerService().RegServer(param)
	raw, err := enPayload(resp, param.ResponseEncryptSecret)
	if err != nil {
		logrus.Warn("报文无法加密")
		return
	}
	topicID := strings.TrimLeft(message.Topic(), regSub[:len(regSub)-1])
	if !c.Publish(regPub+topicID, 0, false, raw).WaitTimeout(waitTimeout) {
		logrus.Warn("mqtt服务器响应失败")
	}
}

// runInitCallBack 初始化回调
func runInitCallBack(c mqtt.Client, message mqtt.Message) {
	info := getAuthInfo(message.Topic())
	if info == nil {
		return
	}
	data := req.ServerInit{}
	err := dePayload(message.Payload(), info.APISecret, &data)
	if err != nil {
		return
	}

	err = SendTelegrafConfig(info.ServerID)
	if err != nil {
		logrus.Warn("telegraf配置发送失败")
	}
}

// SendDeployTask 启动部署
func SendDeployTask(UID, deployID int, serverIds ...int) error {
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
		mqttClient.Publish(deployTopic.Pub(si.MqttUser), 0, false, payloadJson)
	}
	return nil
}

// SendTelegrafConfig 发送telegraf监控配置
func SendTelegrafConfig(serverID int) error {
	defaultConfig := models.GetSysConfigModels().GetTelegraf()
	serverModel := models.GetServerModels()
	serInfo, err := serverModel.Get(serverID)
	if err != nil {
		return err
	}
	selfConfig := serverModel.GetTelegrafConfig(serverID)
	if selfConfig.TomlConfig == "" {
		selfConfig.TomlConfig = defaultConfig
	}
	payload, err := enPayload(resp.TelegrafConfig{
		TomlConfig: selfConfig.TomlConfig,
		IsOpen:     selfConfig.IsOpen != 2,
	}, serInfo.APISecret)
	if err != nil {
		return err
	}
	token := mqttClient.Publish(telegrafConfigTopic.Pub(serInfo.MqttUser), 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// SendTelegrafConfig 发送telegraf监控配置
func UpdateClient(serverID int) error {
	serverModel := models.GetServerModels()
	serInfo, err := serverModel.Get(serverID)
	if err != nil {
		return err
	}
	newSource, err := service.GetServerService().GetClientSource(serInfo.Os, serInfo.Arch)
	if err != nil {
		return err
	}
	if serInfo.ClientVersion == newSource.Version {
		return nil
	}
	payload, err := enPayload(newSource, serInfo.APISecret)
	if err != nil {
		return err
	}
	token := mqttClient.Publish(updateTopic.Pub(serInfo.MqttUser), 0, false, payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// deployCallBack 部署回调
func deployCallBack(c mqtt.Client, message mqtt.Message) {
	info := getAuthInfo(message.Topic())
	if info == nil {
		return
	}
	payload := &req.DeployCallBack{}
	err := dePayload(message.Payload(), info.APISecret, payload)
	if err != nil {
		return
	}
	err = models.GetDeployModels().UpdateServerDeployVersion(payload.DeployID, info.ServerID, payload.Version)
	if err != nil {
		return
	}
}

// metricCallBack 监控指标
func metricCallBack(c mqtt.Client, message mqtt.Message) {
	info := getAuthInfo(message.Topic())
	if info == nil {
		return
	}
	rawData := req.ClientMetricParam{}
	err := dePayload(message.Payload(), info.APISecret, &rawData)
	if err != nil || len(rawData.Metrics) == 0 {
		return
	}
	mds := make([]*write.Point, len(rawData.Metrics))
	for k, v := range rawData.Metrics {
		v.Tags["SERVER_ID"] = strconv.Itoa(info.ServerID)
		v.Tags["SERVER_UID"] = strconv.Itoa(info.UID)
		mds[k] = write.NewPoint(v.Name, v.Tags, v.Fields, time.Unix(v.Timestamp, 0))
	}
	err = models.InfluxDB.Write("client_metric", mds...)
	if err != nil {
		logrus.Warnf("%d号服务器监控指标记录失败", info.ServerID)
	}
}

// logCallBack 日志回调
func logCallBack(c mqtt.Client, message mqtt.Message) {
	info := getAuthInfo(message.Topic())
	if info == nil {
		return
	}
	rawData := make(map[string]interface{})
	err := dePayload(message.Payload(), info.APISecret, &rawData)
	if err != nil {
		return
	}
	mName := "DEFAULT"
	mTags := make(map[string]string)
	mFields := make(map[string]interface{})
	for k, v := range rawData {
		vData := fmt.Sprintf("%v", v)
		switch k {
		case "__MODEL__": // 主题模块
			mName = vData
		case "__MESSAGE__": // 消息体
			mFields["Message"] = vData
		case "__LOG_LEVEL__": // 消息级别
			mTags["LOG_LEVEL"] = vData
		default:
			mTags[k] = vData
		}
	}
	mTags["SERVER_ID"] = strconv.Itoa(info.ServerID)
	mTags["SERVER_UID"] = strconv.Itoa(info.UID)
	mTags["SERVER_DEVICE"] = info.DeviceID
	err = models.InfluxDB.Write("client_log", write.NewPoint(mName, mTags, mFields, time.Now()))
	if err != nil {
		logrus.Warnf("%d号服务器日志记录失败", info.ServerID)
	}
}
