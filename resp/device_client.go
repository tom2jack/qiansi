package resp

type RegServer struct {
	ID        int    `json:"id"`
	Uid       int    `json:"uid"`
	DeviceId  string `json:"device_id"`
	ApiSecret string `json:"api_secret"`
	ErrMsg    string `json:"err_msg"`
}
