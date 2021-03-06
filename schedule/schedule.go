package schedule

import (
	"github.com/jakecoffman/cron"
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common/utils"
	"github.com/zhi-miao/qiansi/models"
	"strconv"
	"time"
)

type task struct {
	// 定时任务句柄
	Cron *cron.Cron
	// 任务队列
	Queue chan models.Schedule
	// 任务执行器数量
	MaxActive int
	// 任务错误重试次数
	TryNum int8
}

type jobResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

type job struct {
	Schedule models.Schedule
}

func (j *job) Run() {
	// 判断是否自动销毁
	if j.Schedule.Remain == 0 || j.Schedule.Remain == 1 {
		Task.Remove(&j.Schedule)
	}
	if j.Schedule.Remain == 0 {
		return
	}
	if j.Schedule.Remain > 0 {
		j.Schedule.Remain--
	}
	Task.Queue <- j.Schedule
}

var Task task

func Start() {
	Task = task{
		Cron:      cron.New(),
		Queue:     make(chan models.Schedule, 50000),
		MaxActive: 200,
		TryNum:    1,
	}
	// 启动任务执行器
	jobRun()
	logrus.Info("开始初始化定时任务")
	scheduleModel := new(models.Schedule)
	taskNum := 0
	lastId := 0
	for {
		taskList, count := scheduleModel.ExportList(lastId, 1)
		if count == 0 {
			break
		}
		for _, item := range taskList {
			if item.Remain == 0 {
				taskNum--
				continue
			}
			err := Task.Add(item)
			if err != nil {
				taskNum--
				logrus.Warn("添加任务到调度器失败", err)
			}
		}
		taskNum += count
		lastId = taskList[count-1].Id
	}
	// 启动任务调度器
	Task.Cron.Start()
	logrus.Infof("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// Add 添加任务
func (t *task) Add(m models.Schedule) error {
	err := utils.PanicToError(func() {
		cronName := strconv.Itoa(m.Id)
		job := &job{Schedule: m}
		Task.Cron.AddJob(m.Crontab, job, cronName)
	})
	return err
}

// 直接运行任务
func (t *task) Run(m *models.Schedule) bool {
	Task.Queue <- *m
	return true
}

// Remove 移除任务
func (t *task) Remove(m *models.Schedule) {
	Task.Cron.RemoveJob(strconv.Itoa(m.Id))
}

// 任务执行器
func jobRun() {
	for i := 0; i < Task.MaxActive; i++ {
		go func(idle int) {
			var m models.Schedule
			for {
				m = <-Task.Queue
				logrus.Infof("Task %d Runing", m.Id)
				beforeExecJob(&m)
				taskResult := execJob(&m)
				logrus.Infof("Task %d End, Result: %#v", m.Id, taskResult)
			}
		}(i)
	}
}

// NextTime 获取下一次执行时间
func (t *task) NextTime(cronTable string) time.Time {
	return cron.Parse(cronTable).Next(time.Now().Local())
}

// 任务前置执行
func beforeExecJob(m *models.Schedule) {
	m.PrevTime = time.Now().Local()
	m.NextTime = Task.NextTime(m.Crontab)
	// 数据库回调
	m.RunCallBack()
}

// 执行具体任务
func execJob(m *models.Schedule) jobResult {
	handler := createHandler(m)
	if handler == nil {
		return jobResult{}
	}
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("panic#schedule/schedule.go:execJob#", err)
		}
	}()
	var i int8 = 0
	var output string
	var err error
	for i < Task.TryNum {
		output, err = handler.Run(m)
		if err == nil {
			return jobResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < Task.TryNum {
			logrus.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", m.Id, i, output, err.Error())
			// 默认重试间隔时间，每次递增10s
			time.Sleep(time.Duration(i) * time.Second * 10)
		}
	}
	return jobResult{Result: output, Err: err, RetryTimes: i}
}
