package udp_hook

import "qiansi/common/utils"

type Hook001PO struct {
	// 部署存储集合
	Deploy *utils.SafeStringMap
}

type Hook001DTO struct {
	// 部署数据交换
	Deploy string
}

func NewHook001PO() *Hook001PO {
	return &Hook001PO{
		Deploy: utils.NewSafeStringMap(),
	}
}
