package main

import (
	"bytes"
	"io"
	"net"
	"strings"
)

const MAX_PEERS = 10

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		if string(buf[:n]) == "send_peers" {
			// Send peers
		}
	}
}

func runServer() {
	l, err := net.Listen("tcp", ":696969")
	logger.Info("Listening on port 696969")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func getPeers(peer string) []string {
	// Retrieve and parse peers list from peer
	conn, err := net.Dial("tcp", peer)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("send_peers"))
	var buf bytes.Buffer
	_, err := io.Copy(&buf, conn)
	if err != nil {
		panic(err)
	}
	return strings.Split(buf.String(), ",")
}
