package udp_service

import (
	"bytes"
	"encoding/gob"
	"qiansi/common/dto"
	"qiansi/common/utils"
)

var Hook001 *Hook001PO

// 001-轮训任务存储器
type Hook001PO struct {
	// 部署存储集合
	Deploy *utils.SafeStringMap
	// 调度任务
	Scheduld *utils.SafeStringMap
}

func init() {
	Hook001 = &Hook001PO{
		Deploy:   utils.NewSafeStringMap(),
		Scheduld: utils.NewSafeStringMap(),
	}
}

// ClientTaskLoop 轮训， Task结构为 {001{deploy_id}:1}
func Hook_001(request []byte) []byte {
	var network bytes.Buffer
	// 1.创建编码器
	enc := gob.NewEncoder(&network)
	// 2.向编码器中写入数据
	err := enc.Encode(dto.Hook001DTO{
		Deploy:   Hook001.Deploy.GET(string(request)),
		Scheduld: Hook001.Deploy.GET(string(request)),
	})
	if err != nil {
		return nil
	}
	return network.Bytes()
}
