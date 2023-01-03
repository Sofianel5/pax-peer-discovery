package main

import (
	"io"
	"log"
	"net"
	"os"
)

const SockAddr = "/tmp/mp-spdz.sock"

func echoServer(c net.Conn) {
	log.Printf("Client connected [%s]", c.RemoteAddr().Network())
	io.Copy(c, c)
	c.Close()
}

func ipcSend(msg string) {
	conn, err := net.Dial("unix", SockAddr)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()
	log.Printf("Sending to server: %s", msg)
	conn.Write([]byte(msg))
}

func listen() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	for {
		// Accept new connections, dispatching them to echoServer
		// in a goroutine.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go echoServer(conn)
	}
}
