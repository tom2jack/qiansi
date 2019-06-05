package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"log"
	"tools-server/conf"
)

var (
	ZM_Clinet *sdk.Client
)

func LoadAliyunSDK() {
	var err error
	ZM_Clinet, err = sdk.NewClientWithAccessKey(
		conf.App.MustValue("Aliyun", "RegionId"),
		conf.App.MustValue("Aliyun", "AccessKey"),
		conf.App.MustValue("Aliyun", "AccessSecret"),
	)
	if err != nil {
		log.Print(err.Error())
	}
}
