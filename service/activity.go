package service

import (
	"github.com/zhi-miao/qiansi/models"
	"math/rand"
)

// InviterActive 邀请活动
func InviterActive(inviter, invitee int) (err error) {
	member := models.Member{Id: inviter}
	switch rand.Intn(2) {
	case 0:
		err = member.AddDeployNum(1)
	case 1:
		err = member.AddScheduleNum(1)
	default:
	}
	return
}
