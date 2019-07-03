package deploy

import (
	"qiansi/common/zmlog"
	"qiansi/models"
	"qiansi/qiansi-client/request"
	"runtime"
)

func Run() {
	TaskList := []models.Deploy{}
	_ = request.GetDeployTask(&TaskList)
	for _, v := range TaskList {
		switch v.DeployType {
		case 1:
			Git(&v)
		}
	}
}

func LogPush(format string, v ...interface{}) {
	// TODO: 推送至千丝平台
	var fname = "default"
	if pc, _, _, ok := runtime.Caller(1); ok {
		fname = runtime.FuncForPC(pc).Name()
	}
	zmlog.Info("("+fname+") "+format, v...)
}
