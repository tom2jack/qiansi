package net_service

import (
	"qiansi/qiansi-server/net_service/udp_service"
)

func init() {
	go udp_service.Start()
}
