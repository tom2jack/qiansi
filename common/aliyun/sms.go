package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
	"tools-server/conf"
)

type smsResponse struct {
	Code      string
	Message   string
	BizId     string
	RequestId string
}

//SendSms 验证码发送
func SendSmsVerify(phone string, code string) bool {
	log_str := fmt.Sprintf("[阿里短信](%s-%s):", phone, code)
	param, err := json.Marshal(map[string]string{
		"verify": code,
	})
	if err != nil {
		log.Print(log_str, err.Error())
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
	response, err := ZM_Clinet.ProcessCommonRequest(request)
	if err != nil {
		log.Print(log_str, err.Error())
		return false
	}
	//{"Message":"OK","RequestId":"B1EB5CD5-1F93-4983-852D-B225B934CF65","BizId":"440415459725393922^0","Code":"OK"}
	str := response.GetHttpContentString()
	json_data := &smsResponse{}
	err = json.Unmarshal([]byte(str), json_data)
	if !response.IsSuccess() || err != nil || json_data.Code != "OK" {
		log.Print(log_str, str)
		return false
	}
	return true
}
