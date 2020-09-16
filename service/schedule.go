package service

import (
	"github.com/zhi-miao/gutils"
	"github.com/zhi-miao/qiansi/common/errors"
	"github.com/zhi-miao/qiansi/common/req"
	"github.com/zhi-miao/qiansi/common/resp"
	"github.com/zhi-miao/qiansi/models"
	"github.com/zhi-miao/qiansi/schedule"
)

type scheduleService struct{}

func GetScheduleService() *scheduleService {
	return &scheduleService{}
}

// Create 创建调度任务
func (s *scheduleService) Create(uid int, param *req.ScheduleCreateParam) error {
	user, err := models.GetMemberModels().UsereInfo(uid)
	if err != nil {
		return err
	}
	scheduleModel := models.GetScheduleModels()
	scheduleNum, err := scheduleModel.Count(uid)
	if err != nil {
		return err
	}
	if scheduleNum >= user.MaxSchedule {
		return errors.NewApiError("您的调度任务创建数量已达上限")
	}
	po := &models.Schedule{}
	gutils.SuperConvert(param, po)
	po.UId = uid
	err = gutils.PanicToError(func() {
		po.NextTime = schedule.Task.NextTime(po.Crontab)
	})
	if err != nil {
		return errors.New("表达式错误")
	}
	if err := scheduleModel.Create(po); err != nil {
		return errors.New("任务档案建立失败")
	}
	err = schedule.Task.Add(*po)
	if err != nil {
		resp.NewApiResult(-5, "添加任务到调度器失败").Json(c)
		return
	}
	resp.NewApiResult(1).Json(c)
	return nil
}
