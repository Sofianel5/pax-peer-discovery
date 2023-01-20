package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
)

const MAX_PEERS = 10
const CONN_PORT = ":42069"

var peers = []string{}
var runningMpc = false

func lockDarkpool() {
	runningMpc = true
}

func unlockDarkpool() {
	runningMpc = false
}

func handleConnection(conn net.Conn, config *Config, myaddr string) {
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
		} else if string(buf[:n]) == "run_darkpool" {
			if !runningMpc {
				lockDarkpool()
				run2pc("dark_pool_inputs", parseHexAddr(config.BuyAsset)+" "+parseHexAddr(config.SellAsset), myaddr, conn.RemoteAddr().String(), "1")
				unlockDarkpool()
			} else {
				logger.Warn("Cannot run darkpool, MPC is already running")
			}
		}
	}
}

func runServer(config *Config, myaddr string) {
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
		go handleConnection(conn, config, myaddr)
	}
}

func runDarkpool(peer, buyAddr, sellAddr, myaddr string) (output string, success error) {
	if runningMpc {
		logger.Warn("Cannot run darkpool, MPC is already running")
		return "", errors.New("Cannot run darkpool, MPC is already running")
	}
	conn, err := net.Dial("tcp", peer+CONN_PORT)
	if err != nil {
		logger.Warn("Could not connect to peer ", peer, ": ", err)
		return "", err
	} else {
		logger.Info("Connected to peer at ", conn.RemoteAddr().String())
	}
	defer conn.Close()
	conn.Write([]byte("run_darkpool"))
	lockDarkpool()
	defer unlockDarkpool()
	return run2pc("dark_pool_inputs", parseHexAddr(buyAddr)+" "+parseHexAddr(sellAddr), myaddr, peer, "0")
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
