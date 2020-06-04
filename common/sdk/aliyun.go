package sdk

import (
	"encoding/json"
	"fmt"
	"gitee.com/zhimiao/qiansi/common"
	aliyunSDK "github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/sirupsen/logrus"
)

type aliyun struct {
	client *aliyunSDK.Client
}

type aliyunSmsResponse struct {
	Code      string
	Message   string
	BizId     string
	RequestId string
}

func NewAliyunSDK() *aliyun {
	client, err := aliyunSDK.NewClientWithAccessKey(
		common.Config.Aliyun.RegionId,
		common.Config.Aliyun.AccessKey,
		common.Config.Aliyun.AccessSecret,
	)
	if err != nil {
		logrus.Error(err.Error())
	}
	return &aliyun{client: client}
}

//SendSms 验证码发送
func (a *aliyun) SendSmsVerify(phone string, code string) bool {
	log_str := fmt.Sprintf("[阿里短信](%s-%s):", phone, code)
	param, err := json.Marshal(map[string]string{
		"verify": code,
	})
	if err != nil {
		logrus.Warn(log_str + err.Error())
		return false
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = common.Config.Aliyun.SmsConfig.RegionId
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = common.Config.Aliyun.SmsConfig.SignName
	request.QueryParams["TemplateCode"] = common.Config.Aliyun.SmsConfig.TemplateCode
	request.QueryParams["TemplateParam"] = string(param)
	response, err := a.client.ProcessCommonRequest(request)
	if err != nil {
		logrus.Warn(log_str + err.Error())
		return false
	}
	//{"Message":"OK","RequestId":"B1EB5CD5-1F93-4983-852D-B225B934CF65","BizId":"440415459725393922^0","Code":"OK"}
	str := response.GetHttpContentString()
	json_data := &aliyunSmsResponse{}
	err = json.Unmarshal([]byte(str), json_data)
	if !response.IsSuccess() || err != nil || json_data.Code != "OK" {
		logrus.Warn(log_str + str)
		return false
	}
	return true
}
