package udp

import (
	"net"
	"os"
	"qiansi/common/conf"
	"qiansi/common/zmlog"
)

func Start() {
	listen := conf.S.MustValue("server", "udp_listen", ":8081")
	addr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		zmlog.Error("Can't resolve address: %v", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		zmlog.Error("Error listening: %v", err)
	}
	defer conn.Close()
	zmlog.Info("Start UDP Service Listening %s", listen)
	for {
		data := make([]byte, 255)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil || n < 4 {
			continue
		}
		go handleClient(data[:n], remoteAddr, conn)
	}
}

// handleClient 取前3位为任务标识位，后为任务数据，最长200
func handleClient(data []byte, remoteAddr *net.UDPAddr, conn *net.UDPConn) {
	var result []byte
	switch string(data[:3]) {
	// 部署任务
	case "001":
		result = ClientTaskLoop(data[3:])
	}
	if result == nil {
		result = []byte("0")
	}
	conn.WriteToUDP(result, remoteAddr)
}
