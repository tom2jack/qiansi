package resp

import "qiansi/qiansi-server/models"

type UserInfoVO struct {
	models.Member
	Token string
}
