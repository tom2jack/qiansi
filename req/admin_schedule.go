package req

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
