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

func ClientTaskLoop(request []byte) []byte {
	have := Task.GET(string(request))
	if have != "" {
		return []byte(have)
	}
	return nil
}
