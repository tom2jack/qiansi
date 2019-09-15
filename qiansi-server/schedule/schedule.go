package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jakecoffman/cron"
	"io/ioutil"
	"net/http"
	"qiansi/qiansi-server/models"
	"strconv"
	"sync"
	"time"
)

type storageData struct {
	Cron   string
	Url    string
	Pay_id string
	Num    int
}

var (
	Cron *cron.Cron
)

func init() {
	Cron = cron.New()
	Cron.Start()
}

func storage_init() {
	schedule := &models.Schedule{
		ScheduleType: 1,
	}
	data, _ := schedule.List(1, 100)

	for i, j := 0, 1; j < len(data); i, j = i+2, j+2 {
		sd := &storageData{}
		xdata, _ := redis.String(data[j], nil)
		err := json.Unmarshal([]byte(xdata), &sd)
		if err != nil {
			fmt.Println("任务解析失败，跳出", xdata)
			continue
		}
		Cron.AddFunc(sd.Cron, func() {
			RequestJob(sd.Url, sd.Pay_id)
		}, sd.Pay_id)
		fmt.Println("成功添加任务：[" + sd.Cron + "]:" + sd.Url)
	}
	fmt.Println("数据导入完毕")
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
