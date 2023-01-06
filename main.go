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
	wg.Add(1)
	for _, peer := range tryPeers {
		wg.Add(1)
		go func(_peer string) {
			success, resp := getPeers(_peer)
			if success {
				logger.Info("Received peers:", resp)
				logger.Info("Now attempting to run MPC protocol with peer ", _peer)
				if _peer[:2] == "54" {
					run2pc("dark_pool_inputs", "611382286831621467233887798921843936019654057231 917551056842671309452305380979543736893630245704", _peer, 0)
				} else {
					run2pc("dark_pool_inputs", "917551056842671309452305380979543736893630245704 611382286831621467233887798921843936019654057231", _peer, 1)
				}
			} else {
				logger.Warn("Could not connect to ", _peer)
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	logger.Info("done")
}
