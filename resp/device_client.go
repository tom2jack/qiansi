package resp

type RegServer struct {
	ID               int    `json:"id"`
	Uid              int    `json:"uid"`
	DeviceID         string `json:"device_id"`
	ApiSecret        string `json:"api_secret"`
	MqttUserName     string `json:"mqtt_user_name"`
	MqttUserPassword string `json:"mqtt_user_password"`
	ErrMsg           string `json:"err_msg"`
}
