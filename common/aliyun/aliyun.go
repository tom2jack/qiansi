package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"qiansi/common/conf"
	"qiansi/common/zmlog"
)

var (
	ZM_Clinet *sdk.Client
)

func init() {
	var err error
	ZM_Clinet, err = sdk.NewClientWithAccessKey(
		conf.S.MustValue("Aliyun", "RegionId"),
		conf.S.MustValue("Aliyun", "AccessKey"),
		conf.S.MustValue("Aliyun", "AccessSecret"),
	)
	if err != nil {
		zmlog.Error(err.Error())
	}
}
