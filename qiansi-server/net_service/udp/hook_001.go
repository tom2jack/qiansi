package udp

import (
	"bytes"
	"encoding/gob"
	"qiansi/common/models/udp_hook"
)

var (
	Hook001 *udp_hook.Hook001PO
)

func init() {
	Hook001 = udp_hook.NewHook001PO()
}

// ClientTaskLoop 轮训， Task结构为 {001{deploy_id}:1}
func Hook_001(request []byte) []byte {
	var network bytes.Buffer
	// 1.创建编码器
	enc := gob.NewEncoder(&network)
	// 2.向编码器中写入数据
	err := enc.Encode(udp_hook.Hook001DTO{
		Deploy: Hook001.Deploy.GET(string(request)),
	})
	if err != nil {
		return nil
	}
	return network.Bytes()
}
