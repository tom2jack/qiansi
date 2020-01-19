package resp

import "gitee.com/zhimiao/qiansi/qiansi-server/models"

type UserInfoVO struct {
	models.Member
	Token string
}
