package service

import "tools-server/service/udp"

func LoadService() {
	go udp.Start()
}
