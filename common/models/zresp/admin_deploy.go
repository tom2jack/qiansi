package zresp

import (
	"qiansi/common/models"
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
	CreateTime    models.JsonTimeDate
	UpdateTime    models.JsonTimeDate
}
