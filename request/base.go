package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	ZMHTTP *http.Client = &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 15,
		},
	}
)

type (
	ZMResponse http.Response
)

type ApiBody struct {
	Code   int
	Result *json.RawMessage
	Msg    string
}

func GET(link string, param url.Values, data *interface{}) error {
	if param.Encode() != "" {
		link = link + "?" + param.Encode()
	}
	raw, err := ZMHTTP.Get(link)
	if err != nil {
		return err
	}
	defer raw.Body.Close()
	body, _ := ioutil.ReadAll(raw.Body)
	requestBody := string(body)
	if raw.Header.Get("DECRYPT") == "1" {
		requestBody = requestBody
	}
	return nil
}

func POST() {

}

type ApiBodyD struct {
	Data []ApiBodyData
}

type ApiBodyData struct {
	server_id  int
	api_secret string
}
