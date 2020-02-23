package req

type VerifyBySMSParam struct {
	// 手机号
	Phone string `json: "phone"`
	// 图片验证码id
	ImgIdKey string `json: "imgIdKey"`
	// 验证码code
	ImgCode string `json:"imgCode"`
}
