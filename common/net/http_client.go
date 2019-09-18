package net

// http-client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type httpClient struct{}

type HttpClientResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

var HttpClient *httpClient

func init() {
	HttpClient = new(httpClient)
}

func (c *httpClient) Get(url string, timeout int) HttpClientResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.createRequestError(err)
	}

	return c.request(req, timeout)
}

func (c *httpClient) PostParams(url string, params string, timeout int) HttpClientResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return c.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return c.request(req, timeout)
}

func (c *httpClient) PostJson(url string, body string, timeout int) HttpClientResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return c.createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return c.request(req, timeout)
}

func (c *httpClient) request(req *http.Request, timeout int) HttpClientResponseWrapper {
	wrapper := HttpClientResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	c.setRequestHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header

	return wrapper
}

func (c *httpClient) setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "golang/gocron")
}

func (c *httpClient) createRequestError(err error) HttpClientResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return HttpClientResponseWrapper{0, errorMessage, make(http.Header)}
}
