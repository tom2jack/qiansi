package req

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
