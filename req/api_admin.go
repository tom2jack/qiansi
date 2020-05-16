package req

import "time"

// DashboardIndexMetricParam 监控入参
type DashboardIndexMetricParam struct {
	ServerId  int `binding:"required"`
	StartTime time.Time
	EndTime   time.Time
}

// 部署应用操作入参
type DeploySetParam struct {
	ID            int    `json:"id"`             // 应用唯一编号
	UId           int    `json:"uid"`            // 创建用户UID
	Title         string `json:"title"`          // 应用名称
	DeployType    int8   `json:"deploy_type"`    // 部署类型 0-本地 1-git 2-zip 3-docker
	WorkDir       string `json:"work_dir"`       // 工作目录
	BeforeCommand string `json:"before_command"` // 前置命令
	AfterCommand  string `json:"after_command"`  // 后置命令
	NowVersion    int    `json:"now_version"`    // 当前版本
	OpenID        string `json:"open_id"`        // 应用开放编码(用于hook部署)
	// DeployDocker 纸喵部署-docker
	DeployDocker struct {
		DockerImage      string `json:"docker_image"`      // 资源地址(完整路径)
		UserName         string `json:"user_name"`         // 账号
		Password         string `json:"password"`          // 密码
		IsRuning         int8   `json:"is_runing"`         // 是否启动 -1-不启动 1-启动
		ContainerName    string `json:"container_name"`    // 容器名称
		ContainerVolumes string `json:"container_volumes"` // 地址映射 a:b a宿主机 b容器内
		ContainerPorts   string `json:"container_ports"`   // 端口暴露 a:b a内部 b外部
		ContainerEnv     string `json:"container_env"`     // 环境变量注入
	} `json:"deploy_docker"`
	// DeployGit 纸喵部署-git
	DeployGit struct {
		RemoteURL  string `json:"remote_url"`  // 资源地址
		DeployPath string `json:"deploy_path"` // 本地部署地址
		Branch     string `json:"branch"`      // git分支
		DeployKeys string `json:"deploy_keys"` // 部署私钥
		UserName   string `json:"user_name"`   // 账号
		Password   string `json:"password"`    // 密码
	} `json:"deploy_git"`
	// DeployZip 纸喵部署-zip
	DeployZip struct {
		RemoteURL  string `json:"remote_url"`  // 资源地址
		DeployPath string `json:"deploy_path"` // 本地部署地址
		Password   string `json:"password"`    // 密码
	} `json:"deploy_zip"`
	ServerRelation []struct {
		// 服务器ID
		ServerId int
		// 关联操作 true-关联 false-取消
		Relation bool
	} `json:"server_relation"`
}

// DeployParam 部署应用入参
type DeployParam struct {
	// 部署应用ID
	DeployId int
}

// DeployListParam 部署应用列表入参
type DeployListParam struct {
	PageParam
	Title string `json:"Title"`
}

type DeployRunLogParam struct {
	DeployParam
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
	DeployParam
}

type DeployServerParam struct {
	DeployParam
}

type DeployDoParam struct {
	DeployParam
}
type DeployRelationParam struct {
	DeployParam
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
	ServerName string `form:"ServerName" json:"ServerName"`
}

type ServerDelParam struct {
	ServerId int `form:"serverId" json:"serverId"`
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
	Phone string `json:"phone"`
	// 图片验证码id
	ImgIdKey string `json:"imgIdKey"`
	// 验证码code
	ImgCode string `json:"imgCode"`
}
