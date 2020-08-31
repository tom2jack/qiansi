package req

// RegServer 注册服务器
type RegServer struct {
	UID                   int    `json:"uid"`
	DeviceID              string `json:"device_id"`
	ResponseEncryptSecret string `json:"response_encrypt_secret"`
	ClientVersion         string `json:"client_version"`
}

type DeployCallBack struct {
	DeployID int `json:"deploy_id"`
	Version  int `json:"version"`
}
