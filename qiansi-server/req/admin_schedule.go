package req

type ScheduleListParam struct {
	PageParam
	Title string
}

type ScheduleCreateParam struct {
	// 执行命令
	Command string `binding:"require"`
	// 表达式
	Crontab string `binding:"require"`
	// 执行次数-1-无限
	Remain int `binding:"require"`
	// 1-http 2-shell
	ScheduleType int `binding:"require"`
	// 执行超时时间
	Timeout int `binding:"require, min=1"`
	// 标题
	Title string `binding:"require"`
	// 执行服务器ID，云服务器id为0
	ServerId int
}

type ScheduleDelParam struct {
	Id int `binding:"min=1"`
}

type ScheduleDoParam struct {
	Id int `binding:"min=1"`
}
