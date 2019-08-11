package api_req

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
