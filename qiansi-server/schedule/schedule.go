package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jakecoffman/cron"
	"io/ioutil"
	"net/http"
	"qiansi/common/logger"
	"qiansi/qiansi-server/models"
	"strconv"
	"time"
)

type task struct {
	Cron *cron.Cron
}

var Task task

func init() {
	Task = task{Cron: cron.New()}
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
			Task.Add(item)
		}
		taskNum += count
		lastId = taskList[count-1].Id
	}
	logger.Info("定时任务初始化完成, 共%d个定时任务添加到调度器", taskNum)
}

func (t *task) Add(m models.Schedule) {

}

func (t *task) Remove(m models.Schedule) {

}

//RequestJob
func RequestJob(url string, pay_id string) {
	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	url = url + "&pay_id=" + string(pay_id) + "&num=" + strconv.Itoa(num)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("正在执行：\nurl:%s\nbody:%s\nnum:\n", url, body)
}

func remove(pay_id string) {
	Cron.RemoveJob(pay_id)
	fmt.Println(pay_id, "删除成功")
}
