package main

import (
	"github.com/zhi-miao/qiansi/api"
	"github.com/zhi-miao/qiansi/common"
	"github.com/zhi-miao/qiansi/device"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/notifyevent"
	"github.com/zhi-miao/qiansi/schedule"
)

// @title 纸喵 qiansi API
// @version 1.0
// @description 纸喵软件系列之服务端
// @termsOfService http://zhimiao.org

// @contact.name API Support
// @contact.url http://tools.zhimiao.org
// @contact.email mail@xiaoliu.org

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http
// @host localhost:1315
// @basepath
func main() {
	common.Config.Init()
	device.Start()
	models.Start()
	go api.Start()
	go schedule.Start()
	go notifyevent.Start()
	select {}
}
