package request

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
	"tools-client/common"
	"tools-client/models"
)

var (
	ZMHTTP *http.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 15,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}
)

const BASE_URL = "https://localhost:8001"

type ApiBody struct {
	Code int
	Data *json.RawMessage
	Msg  string
}

func request(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, BASE_URL+url, body)
	if err != nil {
		return nil, err
	}
	// 设置header头
	if client_id := common.Cfg.String("zhimiao::clientid"); client_id != "" {
		req.Header.Add("SERVER-ID", client_id)
	}
	resp, err := ZMHTTP.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
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
		raw = []byte(common.DecrptogAES(string(raw), common.Cfg.String("zhimiao::apisecret")))
	}
	// 解析统一包装
	apibody := &ApiBody{}
	if err := json.Unmarshal(raw, apibody); err != nil {
		return nil, err
	}
	raw = []byte(*apibody.Data)
	if apibody.Code < 1 {
		return raw, fmt.Errorf(apibody.Msg)
	}
	return raw, nil
}

// 注册当前客户端
func RegServer(uid string, device string) (models.Server, error) {
	url := fmt.Sprintf("/client/ApiRegServer?uid=%s&device=%s", uid, device)
	raw, err := request("GET", url, nil)
	server := models.Server{}
	if err != nil {
		return server, err
	}
	json.Unmarshal(raw, &server)
	return server, nil
}

func GetDeployTask(task *[]models.Deploy) error {
	raw, err := request("GET", "/client/ApiGetDeployTask", nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, task); err != nil {
		return err
	}
	return nil
}
