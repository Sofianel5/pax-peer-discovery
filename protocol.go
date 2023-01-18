package main

import (
	"bytes"
	"io"
	"net"
	"strings"
)

const MAX_PEERS = 10
const CONN_PORT = ":42069"

var peers = []string{}

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

func getPeers(peer string) (conneted bool, peers []string) {
	// Retrieve and parse peers list from peer
	conn, err := net.Dial("tcp", peer+CONN_PORT)
	if err != nil {
		logger.Warn("Could not connect to peer ", peer, ": ", err)
		return false, nil
	} else {
		logger.Info("Connected to peer at ", conn.RemoteAddr().String())
	}
	defer conn.Close()
	conn.Write([]byte("send_peers"))
	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		logger.Warn("Could not read from peer ", peer, ": ", err)
		return false, nil
	}
	return true, strings.Split(buf.String(), ",")
}

func sendPeers(conn net.Conn) {
	// Send peers list to conn
	logger.Info("Sending peers list to ", conn.RemoteAddr().String())
	conn.Write([]byte(strings.Join(peers, ",")))
	conn.Close()
}
