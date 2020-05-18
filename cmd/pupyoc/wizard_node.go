// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Yocoin15/Yocoin_Sources/accounts/keystore"
	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/log"
)

// deployNode creates a new node configuration based on some user input.
func (w *wizard) deployNode(boot bool) {
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
	infos, err := checkNode(client, w.network, boot)
	if err != nil {
		if boot {
			infos = &nodeInfos{portFull: 30303, peersTotal: 512, peersLight: 256}
		} else {
			infos = &nodeInfos{portFull: 30303, peersTotal: 50, peersLight: 0, gasTarget: 4.7, gasPrice: 18}
		}
	}
	existed := err == nil

	infos.genesis, _ = json.MarshalIndent(w.conf.Genesis, "", "  ")
	infos.network = w.conf.Genesis.Config.ChainId.Int64()

	// Figure out where the user wants to store the persistent data
	fmt.Println()
	if infos.datadir == "" {
		fmt.Printf("Where should data be stored on the remote machine?\n")
		infos.datadir = w.readString()
	} else {
		fmt.Printf("Where should data be stored on the remote machine? (default = %s)\n", infos.datadir)
		infos.datadir = w.readDefaultString(infos.datadir)
	}
	if w.conf.Genesis.Config.Ethash != nil && !boot {
		fmt.Println()
		if infos.ethashdir == "" {
			fmt.Printf("Where should the ethash mining DAGs be stored on the remote machine?\n")
			infos.ethashdir = w.readString()
		} else {
			fmt.Printf("Where should the ethash mining DAGs be stored on the remote machine? (default = %s)\n", infos.ethashdir)
			infos.ethashdir = w.readDefaultString(infos.ethashdir)
		}
	}
	// Figure out which port to listen on
	fmt.Println()
	fmt.Printf("Which TCP/UDP port to listen on? (default = %d)\n", infos.portFull)
	infos.portFull = w.readDefaultInt(infos.portFull)

	// Figure out how many peers to allow (different based on node type)
	fmt.Println()
	fmt.Printf("How many peers to allow connecting? (default = %d)\n", infos.peersTotal)
	infos.peersTotal = w.readDefaultInt(infos.peersTotal)

	// Figure out how many light peers to allow (different based on node type)
	fmt.Println()
	fmt.Printf("How many light peers to allow connecting? (default = %d)\n", infos.peersLight)
	infos.peersLight = w.readDefaultInt(infos.peersLight)

	// Set a proper name to report on the stats page
	fmt.Println()
	if infos.yocstats == "" {
		fmt.Printf("What should the node be called on the stats page?\n")
		infos.yocstats = w.readString() + ":" + w.conf.yocstats
	} else {
		fmt.Printf("What should the node be called on the stats page? (default = %s)\n", infos.yocstats)
		infos.yocstats = w.readDefaultString(infos.yocstats) + ":" + w.conf.yocstats
	}
	// If the node is a miner/signer, load up needed credentials
	if !boot {
		if w.conf.Genesis.Config.Ethash != nil {
			// Ethash based miners only need an etherbase to mine against
			fmt.Println()
			if infos.etherbase == "" {
				fmt.Printf("What address should the miner user?\n")
				for {
					if address := w.readAddress(); address != nil {
						infos.etherbase = address.Hex()
						break
					}
				}
			} else {
				fmt.Printf("What address should the miner user? (default = %s)\n", infos.etherbase)
				infos.etherbase = w.readDefaultAddress(common.HexToAddress(infos.etherbase)).Hex()
			}
		} else if w.conf.Genesis.Config.Clique != nil {
			// If a previous signer was already set, offer to reuse it
			if infos.keyJSON != "" {
				if key, err := keystore.DecryptKey([]byte(infos.keyJSON), infos.keyPass); err != nil {
					infos.keyJSON, infos.keyPass = "", ""
				} else {
					fmt.Println()
					fmt.Printf("Reuse previous (%s) signing account (y/n)? (default = yes)\n", key.Address.Hex())
					if w.readDefaultString("y") != "y" {
						infos.keyJSON, infos.keyPass = "", ""
					}
				}
			}
			// Clique based signers need a keyfile and unlock password, ask if unavailable
			if infos.keyJSON == "" {
				fmt.Println()
				fmt.Println("Please paste the signer's key JSON:")
				infos.keyJSON = w.readJSON()

				fmt.Println()
				fmt.Println("What's the unlock password for the account? (won't be echoed)")
				infos.keyPass = w.readPassword()

				if _, err := keystore.DecryptKey([]byte(infos.keyJSON), infos.keyPass); err != nil {
					log.Error("Failed to decrypt key with given passphrase")
					return
				}
			}
		}
		// Establish the gas dynamics to be enforced by the signer
		fmt.Println()
		fmt.Printf("What gas limit should empty blocks target (MGas)? (default = %0.3f)\n", infos.gasTarget)
		infos.gasTarget = w.readDefaultFloat(infos.gasTarget)

		fmt.Println()
		fmt.Printf("What gas price should the signer require (GWei)? (default = %0.3f)\n", infos.gasPrice)
		infos.gasPrice = w.readDefaultFloat(infos.gasPrice)
	}
	// Try to deploy the full node on the host
	nocache := false
	if existed {
		fmt.Println()
		fmt.Printf("Should the node be built from scratch (y/n)? (default = no)\n")
		nocache = w.readDefaultString("n") != "n"
	}
	if out, err := deployNode(client, w.network, w.conf.bootFull, w.conf.bootLight, infos, nocache); err != nil {
		log.Error("Failed to deploy YOC node container", "err", err)
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
