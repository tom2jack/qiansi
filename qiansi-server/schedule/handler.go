package schedule

import (
	"fmt"
	"gitee.com/zhimiao/qiansi/common/net"
	"gitee.com/zhimiao/qiansi/qiansi-server/models"
	"net/http"
)

type Handler interface {
	Run(m *models.Schedule) (string, error)
}

func createHandler(taskModel *models.Schedule) Handler {
	var handler Handler = nil
	switch taskModel.ScheduleType {
	case 1:
		handler = new(HTTPHandler)
	}
	return handler
}

// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(m *models.Schedule) (result string, err error) {
	if m.Timeout <= 0 || m.Timeout > HttpExecTimeout {
		m.Timeout = HttpExecTimeout
	}
	var resp net.HttpClientResponseWrapper
	resp = net.HttpClient.Get(m.Command, m.Timeout)
	// 返回状态码非200，均为失败
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}
	return resp.Body, err
}
