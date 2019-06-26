package clinet_task_loop

import "tools-server/common/utils"

var (
	Task *utils.SafeMap
	RUN  bool = false
)

func init() {
	if RUN {
		return
	}
	Task = utils.NewSafeMap()
	RUN = true
}

// ClientTaskLoop 轮训， Task结构为 {001{deploy_id}:1}
func ClientTaskLoop(request []byte) []byte {
	have := Task.GET(string(request))
	if have != "" {
		return []byte("1")
	}
	return nil
}
