package schedule

import (
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/zhi-miao/qiansi/models"
)

type Handler interface {
	Run(m *models.Schedule) (string, error)
}

func createHandler(taskModel *models.Schedule) Handler {
	var handler Handler
	switch taskModel.ScheduleType {
	case 1:
		handler = new(httpHandler)
	}
	return handler
}

// httpHandler HTTP任务
type httpHandler struct{}

func (h *httpHandler) Run(m *models.Schedule) (result string, err error) {
	var maxTimeout = 50
	if m.Timeout <= 0 || m.Timeout > maxTimeout {
		m.Timeout = maxTimeout
	}
	getRes, err := resty.New().SetTimeout(time.Duration(m.Timeout) * time.Second).R().Get(m.Command)
	if err != nil {
		return "", err
	}
	if getRes.IsSuccess() {
		return string(getRes.Body()), nil
	}
	return "", errors.New(getRes.Status())
}
