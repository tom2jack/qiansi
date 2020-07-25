package notifyevent

import (
	"github.com/sirupsen/logrus"
	"github.com/zhi-miao/qiansi/common"
	"net"
	"os"
)

// 初始化启动监听
func Start() {
	listen := common.Config.Server.UDPListen
	addr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		logrus.Errorf("Can't resolve address: %v", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logrus.Errorf("Error listening: %v", err)
	}
	defer conn.Close()
	logrus.Infof("Start UDP Service Listening %s", listen)
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
