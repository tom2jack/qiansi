package api_req

// 部署应用操作入参
type DeploySetParam struct {
	// 应用唯一编号
	Id int
	// 应用名称
	Title string
	// 部署类型 0-本地 1-git 2-zip
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

type DeployDelParam struct {
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
