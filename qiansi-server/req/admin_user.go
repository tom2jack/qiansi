package req

type UserSiginParam struct {
	// 手机号
	Phone string `json:"phone"`
	// 密码
	Password string
}

type UserSiginUpParam struct {
	UserSiginParam
	// 短信验证码
	Code string
	// 邀请人
	InviterUid int
}

type UserResetPwdParam struct {
	OldPassword string
	NewPassword string
}
