package resp

import (
	"github.com/zhi-miao/qiansi/models"
	"time"
)

type DashboardIndexMetricVO struct {
	ActiveServerNum int                      `json:"active_server_num"`
	CPURate         []map[string]interface{} `json:"cpu_rate"`
	MenRate         []map[string]interface{} `json:"men_rate"`
}

type DashboardInfoVO struct {
	DeployNum    int `json:"deploy_num"`
	MaxDeploy    int `json:"max_deploy"`
	ScheduleNum  int `json:"schedule_num"`
	MaxSchedule  int `json:"max_schedule"`
	ServerNum    int `json:"server_num"`
	InviteNum    int `json:"invite_num"`
	DeployRunNum int `json:"deploy_run_num"`
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
