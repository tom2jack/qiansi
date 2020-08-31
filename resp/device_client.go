package resp

// RegServer 客户端注册返回
type RegServer struct {
	ID               int    `json:"id"`
	Uid              int    `json:"uid"`
	DeviceID         string `json:"device_id"`
	ApiSecret        string `json:"api_secret"`
	MqttUserName     string `json:"mqtt_user_name"`
	MqttUserPassword string `json:"mqtt_user_password"`
	ErrMsg           string `json:"err_msg"`
}

// TelegrafConfig 配置数据
type TelegrafConfig struct {
	TomlConfig string `json:"toml_config"`
	IsOpen     bool   `json:"is_open"`
}
