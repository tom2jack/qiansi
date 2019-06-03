package udp

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func Start() {
	flag.Parse()
	addr, err := net.ResolveUDPAddr("udp", ":8001")
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer conn.Close()
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			continue
		}
		go handleClient(string(data[:n]), remoteAddr, conn)
	}
}

func handleClient(data string, remoteAddr *net.UDPAddr, conn *net.UDPConn) {
	fmt.Println("收到信息:", remoteAddr)
	fmt.Println(data)
	time.Sleep(5 * time.Second)
	conn.WriteToUDP([]byte(data), remoteAddr)
}
