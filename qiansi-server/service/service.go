package service

import (
	udp2 "qiansi/qiansi-server/service/udp"
)

func LoadService() {
	go udp2.Start()
}
