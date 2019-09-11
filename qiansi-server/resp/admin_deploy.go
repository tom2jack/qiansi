package resp

import "time"

type DeployServerVO struct {
	// CreateTime    JsonTimeDate
	DeviceId     string
	Domain       string
	Id           int
	ServerName   string
	ServerRuleId int
	ServerStatus int
	Uid          int
	// UpdateTime    JsonTimeDate
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

type DeployLogVO struct {
	ClientIp      string
	Content       string
	DeployId      int
	DeployVersion int
	DeviceId      string
	Id            int64
	ServerId      int
	Uid           int
	CreateTime    time.Time
}
