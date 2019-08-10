package api_req

type UserSiginParam struct {
	Phone    string
	Password string
}

type UserSiginUpParam struct {
	UserSiginParam
	Code       string
	InviterUid int
}

type UserResetPwdParam struct {
	OldPassword string
	NewPassword string
}
