// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Yocoin15/Yocoin_Sources/log"
)

// deployWallet creates a new web wallet based on some user input.
func (w *wizard) deployWallet() {
	// Do some sanity check before the user wastes time on input
	if w.conf.Genesis == nil {
		log.Error("No genesis block configured")
		return
	}
	if w.conf.yocstats == "" {
		log.Error("No yocstats server configured")
		return
	}
	// Select the server to interact with
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

	// Retrieve any active node configurations from the server
	infos, err := checkWallet(client, w.network)
	if err != nil {
		infos = &walletInfos{
			nodePort: 30303, rpcPort: 8545, webPort: 80, webHost: client.server,
		}
	}
	existed := err == nil

	infos.genesis, _ = json.MarshalIndent(w.conf.Genesis, "", "  ")
	infos.network = w.conf.Genesis.Config.ChainId.Int64()

	// Figure out which port to listen on
	fmt.Println()
	fmt.Printf("Which port should the wallet listen on? (default = %d)\n", infos.webPort)
	infos.webPort = w.readDefaultInt(infos.webPort)

	// Figure which virtual-host to deploy yocstats on
	if infos.webHost, err = w.ensureVirtualHost(client, infos.webPort, infos.webHost); err != nil {
		log.Error("Failed to decide on wallet host", "err", err)
		return
	}
	// Figure out where the user wants to store the persistent data
	fmt.Println()
	if infos.datadir == "" {
		fmt.Printf("Where should data be stored on the remote machine?\n")
		infos.datadir = w.readString()
	} else {
		fmt.Printf("Where should data be stored on the remote machine? (default = %s)\n", infos.datadir)
		infos.datadir = w.readDefaultString(infos.datadir)
	}
	// Figure out which port to listen on
	fmt.Println()
	fmt.Printf("Which TCP/UDP port should the backing node listen on? (default = %d)\n", infos.nodePort)
	infos.nodePort = w.readDefaultInt(infos.nodePort)

	fmt.Println()
	fmt.Printf("Which port should the backing RPC API listen on? (default = %d)\n", infos.rpcPort)
	infos.rpcPort = w.readDefaultInt(infos.rpcPort)

	// Set a proper name to report on the stats page
	fmt.Println()
	if infos.yocstats == "" {
		fmt.Printf("What should the wallet be called on the stats page?\n")
		infos.yocstats = w.readString() + ":" + w.conf.yocstats
	} else {
		fmt.Printf("What should the wallet be called on the stats page? (default = %s)\n", infos.yocstats)
		infos.yocstats = w.readDefaultString(infos.yocstats) + ":" + w.conf.yocstats
	}
	// Try to deploy the wallet on the host
	nocache := false
	if existed {
		fmt.Println()
		fmt.Printf("Should the wallet be built from scratch (y/n)? (default = no)\n")
		nocache = w.readDefaultString("n") != "n"
	}
	if out, err := deployWallet(client, w.network, w.conf.bootFull, infos, nocache); err != nil {
		log.Error("Failed to deploy wallet container", "err", err)
		if len(out) > 0 {
			fmt.Printf("%s\n", out)
		}
		return
	}
	// All ok, run a network scan to pick any changes up
	log.Info("Waiting for node to finish booting")
	time.Sleep(3 * time.Second)

	w.networkStats()
}
