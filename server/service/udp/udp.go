package udp

import (
	"log"
	"net"
	"os"
	"tools-server/conf"
	clinet_task_loop2 "tools-server/server/service/udp/clinet_task_loop"
)

func Start() {
	listen := conf.S.MustValue("server", "udp_listen", ":8081")
	addr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		log.Print("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Print("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Printf("Start UDP Service Listening %s", listen)
	for {
		data := make([]byte, 200)
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
	case "001":
		result = clinet_task_loop2.ClientTaskLoop(data[3:])
	}
	if result == nil {
		result = []byte("0")
	}
	conn.WriteToUDP(result, remoteAddr)
}
