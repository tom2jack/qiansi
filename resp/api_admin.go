package resp

import (
	"gitee.com/zhimiao/qiansi/models"
	"time"
)

type DashboardIndexMetricVO struct {
	ActiveServerNum int
	CPURate         []map[string]interface{}
	MenRate         []map[string]interface{}
}

type DashboardInfoVO struct {
	DeployNum    int
	MaxDeploy    int
	ScheduleNum  int
	MaxSchedule  int
	ServerNum    int
	InviteNum    int
	DeployRunNum int
}

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

type ServerVO struct {
	CreateTime   time.Time
	DeviceId     string
	Domain       string
	Id           int
	ServerName   string
	ServerRuleId int
	ServerStatus int
	UpdateTime   time.Time
}

type UserInfoVO struct {
	models.Member
	Token string
}
