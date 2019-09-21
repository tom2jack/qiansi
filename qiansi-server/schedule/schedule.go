package schedule

import (
	"github.com/jakecoffman/cron"
	"qiansi/common/logger"
	"qiansi/common/utils"
	"qiansi/qiansi-server/models"
	"strconv"
	"time"
)

type task struct {
	Cron     *cron.Cron
	RunQueue *queue
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// 并发队列
type queue struct{ queue chan struct{} }

func (cq *queue) Add()  { cq.queue <- struct{}{} }
func (cq *queue) Done() { <-cq.queue }

var Task task

func init() {
	Task = task{
		Cron:     cron.New(),
		RunQueue: &queue{queue: make(chan struct{}, 100)},
	}

	Task.Cron.Start()
	logger.Info("开始初始化定时任务")
	scheduleModel := new(models.Schedule)
	taskNum := 0
	lastId := 0
	for {
		taskList, count := scheduleModel.ExportList(lastId, 1)
		if count == 0 {
			break
		}
		for _, item := range taskList {
			Task.Add(&item)
		}
		taskNum += count
		lastId = taskList[count-1].Id
	}
	logger.Info("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

// Add 添加任务
func (t *task) Add(m *models.Schedule) {
	taskFunc := createJob(*m)
	if taskFunc == nil {
		logger.Warn("创建任务处理Job失败,不支持的任务协议#", m.ScheduleType)
		return
	}
	cronName := strconv.Itoa(m.Id)
	err := utils.PanicToError(func() {
		Task.Cron.AddFunc(m.Crontab, taskFunc, cronName)
	})
	if err != nil {
		logger.Warn("添加任务到调度器失败#", err)
	}
}

// 直接运行任务
func (t *task) Run(m *models.Schedule) {
	go createJob(*m)()
}

// Remove 移除任务
func (t *task) Remove(m *models.Schedule) {
	Task.Cron.RemoveJob(strconv.Itoa(m.Id))
}

func createJob(m models.Schedule) cron.FuncJob {
	handler := createHandler(&m)
	if handler == nil {
		return nil
	}
	taskFunc := func() {
		Task.RunQueue.Add()
		defer Task.RunQueue.Done()
		// 前置操作
		beforeExecJob(&m)
		logger.Info("开始执行任务#%s#命令-%s", m.Title, m.Command)
		taskResult := execJob(handler, &m)
		logger.Info("任务完成#%s#命令-%s \n 结果: %#v", m.Title, m.Command, taskResult)
	}
	return taskFunc
}

func beforeExecJob(m *models.Schedule) {
	m.PrevTime = time.Now().Local()
	m.NextTime = cron.Parse(m.Crontab).Next(time.Now().Local())
	// 数据库回调
	m.RunCallBack()
	// 判断是否自动销毁
	if m.Remain == 0 {
		Task.Remove(m)
	}
}

// 执行具体任务
func execJob(handler Handler, m *models.Schedule) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("panic#schedule/schedule.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		output, err = handler.Run(m)
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			logger.Warn("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", m.Id, i, output, err.Error())
			// 默认重试间隔时间，每次递增1分钟
			time.Sleep(time.Duration(i) * time.Minute)
		}
	}
	return TaskResult{Result: output, Err: err, RetryTimes: i}
}
