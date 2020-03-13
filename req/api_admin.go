package req

import "time"

// 部署应用操作入参
type DeploySetParam struct {
	// 应用唯一编号
	Id int
	// 应用名称
	Title string
	// 部署类型 1-git 2-zip 3-other
	DeployType int
	// 资源地址
	RemoteUrl string
	// git分支
	Branch string
	// 本地部署地址
	LocalPath string
	// 后置命令
	AfterCommand string
	// 前置命令
	BeforeCommand string
	// 部署私钥
	DeployKeys string
}

type DeployBase struct {
	// 部署应用ID
	DeployId int
}

type DeployListParam struct {
	PageParam
	Title string
}

type DeployRunLogParam struct {
	DeployBase
	// 服务器ID
	ServerId int
	Version  int
}

// DeployLogParam 部署日志接口入参
type DeployLogParam struct {
	PageParam
	// 开始时间
	StartTime time.Time
	// 结束时间
	EndTime time.Time
	// 应用ID
	DeployId int
	// 版本
	DeployVersion int
	// 服务器
	ServerId int
}

type DeployDelParam struct {
	DeployBase
}

type DeployServerParam struct {
	DeployBase
}

type DeployDoParam struct {
	DeployBase
}
type DeployRelationParam struct {
	DeployBase
	// 服务器ID
	ServerId int
	// 关联操作 true-关联 false-取消
	Relation bool
}

type ScheduleListParam struct {
	PageParam
	Title string
}

type ScheduleCreateParam struct {
	// 执行命令
	Command string `binding:"required"`
	// 表达式
	Crontab string `binding:"required"`
	// 执行次数-1-无限
	Remain int `binding:"required"`
	// 1-http 2-shell
	ScheduleType int `binding:"required"`
	// 执行超时时间
	Timeout int `binding:"required,min=1"`
	// 标题
	Title string `binding:"required"`
	// 执行服务器ID，云服务器id为0
	ServerId int
}

type ScheduleDelParam struct {
	Id int `binding:"min=1"`
}

type ScheduleDoParam struct {
	Id int `binding:"min=1"`
}

type ServerListParam struct {
	PageParam
	ServerName string
}

type ServerDelParam struct {
	ServerId int
}

type ServerSetParam struct {
	Id int
	// 服务器备注名
	ServerName string
	// 服务器规则id
	ServerRuleId int
	// 服务器状态 -1-失效 0-待认领 1-已分配通信密钥 2-已绑定
	ServerStatus int
}

type UserSiginParam struct {
	// 手机号
	Phone string `json:"phone"`
	// 密码
	Password string
}

type UserSiginUpParam struct {
	UserSiginParam
	// 短信验证码
	Code string
	// 邀请人
	InviterUid int
}

type UserResetPwdParam struct {
	OldPassword string
	NewPassword string
}

type VerifyBySMSParam struct {
	// 手机号
	Phone string `json: "phone"`
	// 图片验证码id
	ImgIdKey string `json: "imgIdKey"`
	// 验证码code
	ImgCode string `json:"imgCode"`
}
