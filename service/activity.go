package service

import (
	"math/rand"
	"time"

	"github.com/zhi-miao/qiansi/models"
)

type activityService struct{}

func GetActivityService() *activityService {
	return &activityService{}
}

// Inviter 邀请活动
func (s *activityService) Inviter(inviter, invitee int) (err error) {
	memberModel := models.GetMemberModels()
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(2) {
	case 0:
		err = memberModel.IncrMaxDeployNum(inviter, 1)
	case 1:
		err = memberModel.IncrMaxScheduleNum(inviter, 1)
	default:
	}
	return
}
