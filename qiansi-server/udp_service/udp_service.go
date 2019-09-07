package udp_service

import (
	"net"
	"os"
	"qiansi/common/conf"
	"qiansi/common/logger"
)

// 初始化启动监听
func init() { go Start() }
func Start() {
	listen := conf.S.MustValue("server", "udp_listen", ":8081")
	addr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		logger.Error("Can't resolve address: %v", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Error("Error listening: %v", err)
	}
	defer conn.Close()
	logger.Info("Start UDP Service Listening %s", listen)
	for {
		data := make([]byte, 255)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil || n < 4 {
			continue
		}
		go func(data []byte, remoteAddr *net.UDPAddr, conn *net.UDPConn) {
			var result []byte
			switch string(data[:3]) {
			// 部署任务
			case "001":
				result = Hook_001(data[3:])
			}
			if result == nil {
				result = []byte("0")
			}
			conn.WriteToUDP(result, remoteAddr)
		}(data[:n], remoteAddr, conn)
	}
}
