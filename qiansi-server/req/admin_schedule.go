package req

type ScheduleListParam struct {
	PageParam
	Title string
}

type ScheduleCreateParam struct {
	Command      string
	Crontab      string
	Remain       int
	ScheduleType int
	Timeout      int
	Title        string
}

type ScheduleDelParam struct {
	Id int `binding:"min=1"`
}

type ScheduleDoParam struct {
	Id int `binding:"min=1"`
}
