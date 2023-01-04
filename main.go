package main

import (
	"github.com/ipfs/go-log/v2"
)

var logger = log.Logger("darkpool")

func main() {
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel("darkpool", "info")
	logger.Info("Hello World, starting node...")
	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}
	var tryPeers = findPeers(config)
	logger.Info("Found public peers:", tryPeers)
	// ipcSend("/try/10.0.0.1")
	go runServer()
	for _, peer := range tryPeers {
		// go getPeers(peer)
		go func(_peer string) {
			resp := getPeers(_peer)
			logger.Info("Received peers:", resp)
		}(peer)
	}
}
