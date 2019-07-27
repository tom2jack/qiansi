package udp_service

import (
	"bytes"
	"encoding/gob"
	"log"
	"qiansi/common/models/udp_hook"
	"testing"
)

func TestHook_001(t *testing.T) {
	Hook001.Deploy.SET("1111", "11")
	r := Hook_001([]byte("1111"))
	log.Println(len(r))

	network := bytes.NewBuffer(r)
	// 3.创建解码器
	dec := gob.NewDecoder(network)
	var d = &udp_hook.Hook001DTO{}
	// 4.解码
	dec.Decode(&d)
	log.Printf("%#v", d)
}
