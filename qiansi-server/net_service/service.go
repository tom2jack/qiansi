package net_service

import (
	udp2 "qiansi/qiansi-server/net_service/udp"
)

func LoadService() {
	go udp2.Start()
}
