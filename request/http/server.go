package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"tools-client/models"
)

func RegServer(uid string, device string) (models.Server, error) {
	url := fmt.Sprintf(BASE_URL+"/client/ApiRegServer?uid=%s&device=%s", uid, device)
	fmt.Println(url)
	server := models.Server{}
	raw, err := ZMHTTP.Get(url)
	if err != nil {
		return server, err
	}
	defer raw.Body.Close()
	if raw.StatusCode != 200 {
		return server, fmt.Errorf("非法请求")
	}
	body, _ := ioutil.ReadAll(raw.Body)
	type ApiBody struct {
		Code int
		Data models.Server
		Msg  string
	}
	apibody := &ApiBody{}
	json.Unmarshal(body, apibody)
	if apibody.Code != 1 {
		return server, fmt.Errorf(apibody.Msg)
	}
	return apibody.Data, nil
}
