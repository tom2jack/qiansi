package req

// ClientMetricParam 客户端指标参数
type ClientMetricParam struct {
	Metrics []struct {
		Fields    map[string]interface{} `json:"fields"`
		Name      string                 `json:"name"`
		Tags      map[string]string      `json:"tags"`
		Timestamp int64                  `json:"timestamp"`
	} `json:"metrics"`
}

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

// 遥测请求
type TelesignalParam struct {
	CheckOnline bool // 是否在线
}
