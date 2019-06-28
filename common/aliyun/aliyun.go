package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"log"
	"qiansi/conf"
)

var (
	ZM_Clinet *sdk.Client
)

func LoadAliyunSDK() {
	var err error
	ZM_Clinet, err = sdk.NewClientWithAccessKey(
		conf.S.MustValue("Aliyun", "RegionId"),
		conf.S.MustValue("Aliyun", "AccessKey"),
		conf.S.MustValue("Aliyun", "AccessSecret"),
	)
	if err != nil {
		log.Print(err.Error())
	}
}
