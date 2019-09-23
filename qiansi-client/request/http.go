package request

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"qiansi/common/conf"
	"qiansi/common/dto"
	"qiansi/common/logger"
	"qiansi/common/utils"
	"strconv"
	"strings"
	"time"
)

var (
	ZMHTTP *http.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 15,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
)

const BASE_URL = ""

type ApiBody struct {
	Code int
	Data *json.RawMessage
	Msg  string
}

func request(method string, url string, body io.Reader) ([]byte, error) {
	var real_metod = method
	if real_metod == "RAW" {
		real_metod = "POST"
	}
	req, err := http.NewRequest(real_metod, conf.API_URL+url, body)
	if err != nil {
		return nil, err
	}
	// 设置header头
	if client_id := conf.C.MustValue("zhimiao", "clientid"); client_id != "" {
		req.Header.Add("SERVER-ID", client_id)
	}
	if device := conf.C.MustValue("zhimiao", "device"); device != "" {
		req.Header.Add("SERVER-DEVICE", device)
	}
	if uid := conf.C.MustValue("zhimiao", "uid"); uid != "" {
		req.Header.Add("SERVER-UID", uid)
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := ZMHTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 请求判断
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("非法请求")
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 解密判断
	if resp.Header.Get("ZHIMIAO-Encypt") == "1" {
		raw = []byte(utils.DecrptogAES(string(raw), conf.C.MustValue("zhimiao", "apisecret")))
	}
	// 解析统一包装
	apibody := &ApiBody{}
	if err := json.Unmarshal(raw, apibody); err != nil {
		return nil, err
	}
	if apibody.Code < 1 {
		return raw, fmt.Errorf(apibody.Msg)
	}
	// 防止没有返回数据的情况下发生painc错误
	if apibody.Data != nil {
		raw = []byte(*apibody.Data)
	} else {
		raw = nil
	}
	logger.Info("[发送请求]:(%s)%s\n[返回结果]:%s", method, conf.API_URL+url, string(raw))
	return raw, nil
}

// 注册当前客户端
func RegServer(uid string, device string) (dto.ServerDTO, error) {
	url := fmt.Sprintf("/client/ApiRegServer?uid=%s&device=%s", uid, device)
	raw, err := request("GET", url, nil)
	server := dto.ServerDTO{}
	if err != nil {
		return server, err
	}
	json.Unmarshal(raw, &server)
	return server, nil
}

func GetDeployTask(task *[]dto.DeployDTO) error {
	raw, err := request("GET", "/client/ApiGetDeployTask", nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, task); err != nil {
		return err
	}
	return nil
}

// 推送日志
func LogPush(deploy *dto.DeployDTO, log string) error {
	body := url.Values{
		"deployId": {strconv.Itoa(deploy.Id)},
		"version":  {strconv.Itoa(deploy.NowVersion)},
		"content":  {log},
	}
	_, err := request("POST", "/client/ApiDeployLog", strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	return nil
}

// 部署任务回调
func DeployNotify(deploy *dto.DeployDTO) {
	url := fmt.Sprintf("/client/ApiDeployNotify?deployId=%d&version=%d", deploy.Id, deploy.NowVersion)
	_, err := request("GET", url, nil)
	if err != nil {
		return
	}
}
