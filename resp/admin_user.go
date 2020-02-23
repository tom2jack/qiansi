package resp

import "gitee.com/zhimiao/qiansi/models"

type UserInfoVO struct {
	models.Member
	Token string
}
