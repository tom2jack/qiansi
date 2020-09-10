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

// ServerInit 客户端上线初始化
type ServerInit struct {
	ClientVersion string `json:"client_version"` // 客户端版本
	OS            string `json:"os"`             // 客户端运行系统
	Arch          string `json:"arch"`           // 硬件架构
}
