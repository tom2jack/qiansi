package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/zhi-miao/qiansi/resp"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
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
	resp := service.RegServer(param)
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
	err := SendTelegrafConfig(info.ServerID)
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
		mqttClient.Publish(fmt.Sprintf(deployPub, si.MqttUser), 0, false, payloadJson)
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
	token := mqttClient.Publish(fmt.Sprintf(telegrafConfigPub, serInfo.MqttUser), 0, false, payload)
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
	err = models.GetServerModels().UpdateDeployVersion(payload.DeployID, info.ServerID, payload.Version)
	if err != nil {
		return
	}
}

// metricCallBack 监控指标
func metricCallBack(c mqtt.Client, message mqtt.Message) {
	raw, err := c.GetRawData()
	if err != nil {
		c.Status(403)
		return
	}
	rawData := req.ClientMetricParam{}
	err = json.Unmarshal(raw, &rawData)
	if err != nil {
		fmt.Print(err.Error())
		c.Status(400)
		return
	}
	if len(rawData.Metrics) == 0 {
		resp.NewApiResult(1).Json(c)
		return
	}
	serverId := strconv.Itoa(c.GetInt("SERVER-ID"))
	serverUid := strconv.Itoa(c.GetInt("SERVER-UID"))
	mds := make([]*write.Point, len(rawData.Metrics))
	for k, v := range rawData.Metrics {
		v.Tags["SERVER_ID"] = serverId
		v.Tags["SERVER_UID"] = serverUid
		mds[k] = write.NewPoint(v.Name, v.Tags, v.Fields, time.Unix(v.Timestamp, 0))
	}
	models.InfluxDB.Write("client_metric", mds...)
}

// logCallBack 日志回调
func logCallBack(c mqtt.Client, message mqtt.Message) {
	// TODO
	raw, err := c.GetRawData()
	if err != nil {
		c.Status(403)
		return
	}
	rawData := map[string]interface{}{}
	err = json.Unmarshal(raw, &rawData)

	if err != nil {
		c.Status(400)
		return
	}
	mFields := map[string]interface{}{}
	mName := ""
	mTags := map[string]string{
		"SERVER_ID":     strconv.Itoa(c.GetInt("SERVER-ID")),
		"SERVER_UID":    strconv.Itoa(c.GetInt("SERVER-UID")),
		"SERVER_DEVICE": c.GetString("SERVER-DEVICE"),
	}
	for k, v := range rawData {
		vData := fmt.Sprintf("%v", v)
		switch k {
		case "__MODEL__":
			mName = vData
		case "__MESSAGE__":
			mFields["Message"] = vData
		case "__LOG_LEVEL__":
			mTags["LOG_LEVEL"] = vData
		default:
			mTags[k] = vData
		}
	}
	metric := write.NewPoint(mName, mTags, mFields, time.Now())
	models.InfluxDB.Write("client_log", metric)
}
