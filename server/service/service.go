package service

import (
	udp2 "tools-server/server/service/udp"
)

func LoadService() {
	go udp2.Start()
}
