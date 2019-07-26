package udp

import "qiansi/common/utils"

var (
	Task *utils.SafeMap
)

func init() {
	Task = utils.NewSafeMap()
}

// ClientTaskLoop 轮训， Task结构为 {001{deploy_id}:1}
func ClientTaskLoop(request []byte) []byte {
	have := Task.GET(string(request))
	if have != "" {
		return []byte(have)
	}
	return nil
}
