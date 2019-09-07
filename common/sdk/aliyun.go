package sdk

import (
	"encoding/json"
	"fmt"
	aliyunSDK "github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"qiansi/common/conf"
	"qiansi/common/logger"
)

var (
	Aliyun *aliyun
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

func init() {
	var err error
	Aliyun.client, err = aliyunSDK.NewClientWithAccessKey(
		conf.S.MustValue("Aliyun", "RegionId"),
		conf.S.MustValue("Aliyun", "AccessKey"),
		conf.S.MustValue("Aliyun", "AccessSecret"),
	)
	if err != nil {
		logger.Error(err.Error())
	}
}

//SendSms 验证码发送
func (a *aliyun) SendSmsVerify(phone string, code string) bool {
	log_str := fmt.Sprintf("[阿里短信](%s-%s):", phone, code)
	param, err := json.Marshal(map[string]string{
		"verify": code,
	})
	if err != nil {
		logger.Warn(log_str + err.Error())
		return false
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = conf.S.MustValue("AliyunSms", "RegionId")
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = conf.S.MustValue("AliyunSms", "SignName")
	request.QueryParams["TemplateCode"] = conf.S.MustValue("AliyunSms", "TemplateCode")
	request.QueryParams["TemplateParam"] = string(param)
	response, err := a.client.ProcessCommonRequest(request)
	if err != nil {
		logger.Warn(log_str + err.Error())
		return false
	}
	//{"Message":"OK","RequestId":"B1EB5CD5-1F93-4983-852D-B225B934CF65","BizId":"440415459725393922^0","Code":"OK"}
	str := response.GetHttpContentString()
	json_data := &aliyunSmsResponse{}
	err = json.Unmarshal([]byte(str), json_data)
	if !response.IsSuccess() || err != nil || json_data.Code != "OK" {
		logger.Warn(log_str + str)
		return false
	}
	return true
}
