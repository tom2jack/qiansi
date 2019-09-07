package resp

import (
	"qiansi/qiansi-server/models"
)

type DeployServerVO struct {
	models.Server
	models.DeployServerRelation
}

type DeployVO struct {
	AfterCommand  string
	BeforeCommand string
	Branch        string
	DeployKeys    string
	DeployType    int
	Id            int
	LocalPath     string
	NowVersion    int
	RemoteUrl     string
	Title         string
	Uid           int
	CreateTime    JsonTimeDate
	UpdateTime    JsonTimeDate
}
