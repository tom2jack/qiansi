package api_resp

import "qiansi/common/models"

type DeployServerVO struct {
	models.Server
	models.DeployServerRelation
}
