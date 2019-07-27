package net_service

import (
	"qiansi/qiansi-server/net_service/udp_service"
)

func LoadService() {
	go udp_service.Start()
}
