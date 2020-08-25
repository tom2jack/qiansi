package req

type RegServer struct {
	UID                   int    `json:"uid"`
	DeviceID              string `json:"device_id"`
	ResponseEncryptSecret string `json:"response_encrypt_secret"`
}
