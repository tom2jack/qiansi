package deploy

import (
	"qiansi/models"
	"qiansi/qiansi-client/request"
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
