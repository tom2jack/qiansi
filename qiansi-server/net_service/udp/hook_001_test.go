package udp

import (
	"bytes"
	"encoding/gob"
	"log"
	"qiansi/common/models/udp_hook"
	"testing"
)

func TestHook_001(t *testing.T) {
	Hook001.Deploy.SET("111", "this is data ooo")
	r := Hook_001([]byte("111"))
	log.Println(string(r))

	network := bytes.NewBuffer(r)
	// 3.创建解码器
	dec := gob.NewDecoder(network)
	var d = &udp_hook.Hook001DTO{}
	// 4.解码
	dec.Decode(&d)
	log.Printf("%#v", d)
}
