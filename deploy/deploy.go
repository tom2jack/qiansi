package deploy

import "fmt"

// 部署信息配置
type DeployConfig struct {
	SourcePath   string
	LocalPath    string
	SourceBranch string
}

func Run() {
	// TODO: 任务查询，分发执行
	fmt.Print("执行查询！")
}
