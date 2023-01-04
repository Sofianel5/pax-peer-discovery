package main

import (
	"github.com/ipfs/go-log/v2"
	"sync"
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
	go runServer()
	var tryPeers = findPeers(config)
	logger.Info("Found public peers:", tryPeers)
	// ipcSend("/try/10.0.0.1")
	var wg sync.WaitGroup
	for _, peer := range tryPeers {
		wg.Add(1)
		go func(_peer string) {
			success, resp := getPeers(_peer)
			if success {
				logger.Info("Received peers:", resp)
			} else {
				logger.Warn("Could not connect to ", _peer)
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	logger.Info("done")
}
