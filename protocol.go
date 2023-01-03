package main

import (
	"net"
)

func checkPeerHello(ip string) bool {
	// send TCP hello to peer

	var conn net.Conn
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		logger.Error("Error connecting to peer:", err)
		return false
	}
	defer conn.Close()
	return true
}
