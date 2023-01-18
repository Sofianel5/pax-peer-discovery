package main

import (
	"github.com/ipfs/go-log/v2"
	"sync"
)

var logger = log.Logger("darkpool")

func main() {
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel("darkpool", "info")
	log.SetLogLevel("net/identify", "error")
	log.SetLogLevel("dht/RtRefreshManager", "error")
	logger.Info("Hello World, starting node...")
	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}
	myaddr := getMyIp()
	go runServer(&config, myaddr)
	var tryPeers = findPeers(config, myaddr)
	logger.Info("Found public peers:", tryPeers)
	// ipcSend("/try/10.0.0.1")
	var wg sync.WaitGroup
	wg.Add(1)
	for _, peer := range tryPeers {
		wg.Add(1)
		go func(_peer string) {
			output, err := runDarkpool(_peer, parseHexAddr(config.BuyAsset), parseHexAddr(config.SellAsset), myaddr)
			if err != nil {
				logger.Warn("Could not run darkpool on peer ", _peer, ": ", err)
			} else {
				logger.Info("Darkpool output: ", output)
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()
	logger.Info("done")
}
