package api

import (
	"fmt"
	"log"
	"net"
)

func GetProcessPort() {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("无法绑定端口：%v", err)
	}
	defer ln.Close()

	address := ln.Addr().(*net.TCPAddr)
	port := address.Port

	fmt.Printf("当前进程绑定的端口是：%d\n", port)
}
