package main

import (
	"bytes"
	"io"
	"net"
	"strings"
)

const MAX_PEERS = 10
const CONN_PORT = ":696969"

var peers = make([]string, MAX_PEERS)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		logger.Info("Received:", string(buf[:n]))
		if string(buf[:n]) == "send_peers" {
			// Send peers
			sendPeers(conn)
		}
	}
}

func runServer() {
	l, err := net.Listen("tcp", CONN_PORT)
	logger.Info("Listening on port", CONN_PORT)
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
	conn, err := net.Dial("tcp", peer+CONN_PORT)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Write([]byte("send_peers"))
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		panic(err)
	}
	return strings.Split(buf.String(), ",")
}

func sendPeers(conn net.Conn) {
	// Send peers list to conn
	conn.Write([]byte(strings.Join(peers, ",")))
}
