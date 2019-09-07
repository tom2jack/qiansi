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

type SafeTask struct {
	sync.RWMutex
	Map map[string]int
}

type listData struct {
	Next string
	Prev string
	Name string
}

type storageData struct {
	Cron   string
	Url    string
	Pay_id string
	Num    int
}

var (
	Cron      = cron.New()
	Task      = newSafeTask()
	RedisPool = newRedisPool()
)

func main() {
	Cron.Start()
	storage_init()
	Server()
}

func newSafeTask() *SafeTask {
	st := new(SafeTask)
	st.Map = make(map[string]int)
	return st
}

func (st *SafeTask) readMap(key string) int {
	st.RLock()
	value := st.Map[key]
	st.RUnlock()
	return value
}

func (st *SafeTask) writeMap(key string, value int) {
	st.Lock()
	st.Map[key] = value
	st.Unlock()
}

func (st *SafeTask) delMap(key string) {
	st.Lock()
	delete(st.Map, key)
	st.Unlock()
}

func newRedisPool() *redis.Pool {
	return models.Redis
}

func storage_init() {
	fmt.Println("正在导入数据...")
	c := RedisPool.Get()
	defer c.Close()
	data, err := redis.Values(c.Do("HGETALL", "QIANSI:schedule:server"))
	if err != nil {
		fmt.Println("redis 数据读取失败:", err.Error())
	}

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

func storage_set(key string, value string) {
	c := RedisPool.Get()
	c.Do("HSET", "QIANSI:schedule:server", key, value)
	c.Close()
}

func storage_del(key string) {
	c := RedisPool.Get()
	c.Do("HDEL", "QIANSI:schedule:server", key)
	c.Close()
}

//Server
func Server() {
	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		// 设置持久化
		data := &storageData{
			Url:    request.Form.Get("url"),
			Cron:   request.Form.Get("cron"),
			Pay_id: request.Form.Get("pay_id"),
			Num:    1,
		}
		str, _ := json.Marshal(data)
		storage_set(request.Form.Get("pay_id"), string(str))
		Cron.AddFunc(request.Form.Get("cron"), func() {
			RequestJob(request.Form.Get("url"), request.Form.Get("pay_id"))
		}, request.Form.Get("pay_id"))
		writer.Write([]byte("succ"))
		fmt.Println("成功添加任务：[" + request.Form.Get("cron") + "]:" + request.Form.Get("url"))
	})
	http.HandleFunc("/del", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		remove(request.Form.Get("pay_id"))
		writer.Write([]byte("succ"))
	})
	http.HandleFunc("/list", func(writer http.ResponseWriter, request *http.Request) {
		list := make([]listData, 0)
		for _, e := range Cron.Entries() {
			list = append(list, listData{
				Next: e.Schedule.Next(e.Next).Format("2006-01-02 15:04:05"),
				Prev: e.Next.Format("2006-01-02 15:04:05"),
				Name: e.Name,
			})
		}

		b, err := json.Marshal(list)
		if err != nil {
			fmt.Println(err.Error())
		}
		writer.Write(b)
	})
	err := http.ListenAndServe(":52004", nil)
	if err != nil {
		fmt.Println("http服务启动失败:" + err.Error())
	}
}

//RequestJob
func RequestJob(url string, pay_id string) {
	num := Task.readMap(pay_id)
	Task.writeMap(pay_id, num+1)
	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	url = url + "&pay_id=" + string(pay_id) + "&num=" + strconv.Itoa(num)
	resp, err := client.Get(url)
	if err != nil {
		if num > 100 {
			remove(pay_id)
		}
		fmt.Println(err.Error())
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("正在执行：\nurl:%s\nbody:%s\nnum:%d\n", url, body, num)
	if string(body) != "2" || num > 300 {
		remove(pay_id)
		return
	}

	fmt.Printf("正在执行：%s-%d\n", body, num)
}

func remove(pay_id string) {
	Cron.RemoveJob(pay_id)
	Task.delMap(pay_id)
	storage_del(pay_id)
	fmt.Println(pay_id, "删除成功")
}
