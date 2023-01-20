package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
	"sync"
)

func filterPeers(addrList []string, myIp string) []string {
	type void struct{}
	var member void
	var publicIps = make(map[string]void)
	// var myIp = getMyIp()
	for _, addr := range addrList {
		if _, exists := publicIps[addr]; !exists && checkIp(addr) && addr != myIp {
			// check if addr in publicIps
			publicIps[addr] = member
		}
	}
	keys := make([]string, len(publicIps))
	i := 0
	for k := range publicIps {
		keys[i] = k
		i++
	}
	return keys
}

func findPeers(config Config, myIp string) []string {
	host, err := libp2p.New(libp2p.ListenAddrs([]multiaddr.Multiaddr(config.ListenAddresses)...))
	if err != nil {
		panic(err)
	} else {
		logger.Info("Host created. We are:", host.ID())
		logger.Info(host.Addrs())
	}
	ctx := context.Background()
	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}
	logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	logger.Info(config.BootstrapPeers)
	for _, peerAddy := range config.BootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddy)
		if peerInfo != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := host.Connect(ctx, *peerInfo); err != nil {
					logger.Warning(err)
				} else {
					// logger.Info("Connection established with bootstrap node:", *peerInfo)
				}
			}()
		}
	}
	wg.Wait()

	logger.Info("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, config.RendezvousString)
	logger.Debug("Successfully announced!")

	logger.Info("Searching for other peers...")
	peersChan, err := routingDiscovery.FindPeers(ctx, config.RendezvousString)
	if err != nil {
		logger.Error("Error finding peers: ", err)
		panic(err)
	} else {
		logger.Info("Found peers!")
	}
	var foundPeers = make([]string, 0)
	for peer := range peersChan {
		if peer.ID == host.ID() {
			continue
		}
		// logger.Info("Found peer:", peer)
		if len(peer.Addrs) == 0 {
			// logger.Warning("No addresses found for peer:", peer)
			continue
		} else {
			peerIps := make([]string, 0)
			for _, addr := range peer.Addrs {
				val, err := addr.ValueForProtocol(multiaddr.P_IP4)
				if err == nil {
					peerIps = append(peerIps, val)
				}
			}
			foundPeers = append(foundPeers, peerIps...)
		}
	}
	logger.Info("Done searching for peers!")
	var filteredPeers = filterPeers(foundPeers, myIp)
	// Try again if no peers found
	if len(filteredPeers) == 0 {
		logger.Warning("No peers found. Trying again...")
		return findPeers(config, myIp)
	}
	return filteredPeers
}
