package api_resp

import "qiansi/common/models"

type UserInfoVO struct {
	models.Member
	Token string
}
