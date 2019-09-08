package resp

import (
	"time"
)

type DeployServerVO struct {
	CreateTime    time.Time
	DeviceId      string
	Domain        string
	Id            int
	ServerName    string
	ServerRuleId  int
	ServerStatus  int
	Uid           int
	UpdateTime    time.Time
	DeployId      int
	DeployVersion int
	ServerId      int
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
