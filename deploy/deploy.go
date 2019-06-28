package deploy

import (
	"tools-client/models"
	"tools-client/request"
)

func Run() {
	TaskList := []models.Deploy{}
	request.GetDeployTask(&TaskList)
	for _, v := range TaskList {
		switch v.DeployType {
		case 1:
			Git(&v)
		}
	}
}
