package resp

import "time"

// 部署服务器
type DeployServerVO struct {
	DeviceId      string
	Domain        string
	Id            int
	ServerName    string
	ServerRuleId  int
	ServerStatus  int
	Uid           int
	DeployId      int
	DeployVersion int
	ServerId      int
}

type DeployRunLogTabVO struct {
	Id            int
	ServerName    string
	DeployVersion int
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
	OpenId        string
	RemoteUrl     string
	Title         string
	Uid           int
	CreateTime    time.Time
	UpdateTime    time.Time
}
