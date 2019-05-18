package main

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"tools-client/install"
)

var (
	NAME_1 = 1
	Cron   = cron.New()
	client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
)

func main() {
	//Cron.Start()
	//Server()
	//config :=
	//deploy.Git(&deploy.DeployConfig{})
	install.Install()
}

func Server() {
	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		Cron.AddFunc(request.Form.Get("cron"), func() {
			requestJob(request.Form.Get("url"))
		}, request.Form.Get("name"))
		println("成功添加任务：[" + request.Form.Get("cron") + "]:" + request.Form.Get("url"))
	})
	http.HandleFunc("/del", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		Cron.RemoveJob(request.Form.Get("name"))
		println("删除成功")
	})
	log.Fatal(http.ListenAndServe(":52004", nil))
}

func requestJob(url string) {
	resp, err := client.Get(url)
	if err != nil {
		println(err.Error())
		return
	}
	// 完成后关闭流
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("正在执行：%s\n", body)
}
