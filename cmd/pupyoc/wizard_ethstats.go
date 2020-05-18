// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package main

import (
	"fmt"
	"sort"

	"github.com/Yocoin15/Yocoin_Sources/log"
)

// deployEthstats queries the user for various input on deploying an yocstats
// monitoring server, after which it executes it.
func (w *wizard) deployEthstats() {
	// Select the server to interact with
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

	// Retrieve any active yocstats configurations from the server
	infos, err := checkEthstats(client, w.network)
	if err != nil {
		infos = &yocstatsInfos{
			port:   80,
			host:   client.server,
			secret: "",
		}
	}
	existed := err == nil

	// Figure out which port to listen on
	fmt.Println()
	fmt.Printf("Which port should yocstats listen on? (default = %d)\n", infos.port)
	infos.port = w.readDefaultInt(infos.port)

	// Figure which virtual-host to deploy yocstats on
	if infos.host, err = w.ensureVirtualHost(client, infos.port, infos.host); err != nil {
		log.Error("Failed to decide on yocstats host", "err", err)
		return
	}
	// Port and proxy settings retrieved, figure out the secret and boot yocstats
	fmt.Println()
	if infos.secret == "" {
		fmt.Printf("What should be the secret password for the API? (must not be empty)\n")
		infos.secret = w.readString()
	} else {
		fmt.Printf("What should be the secret password for the API? (default = %s)\n", infos.secret)
		infos.secret = w.readDefaultString(infos.secret)
	}
	// Gather any blacklists to ban from reporting
	if existed {
		fmt.Println()
		fmt.Printf("Keep existing IP %v blacklist (y/n)? (default = yes)\n", infos.banned)
		if w.readDefaultString("y") != "y" {
			// The user might want to clear the entire list, although generally probably not
			fmt.Println()
			fmt.Printf("Clear out blacklist and start over (y/n)? (default = no)\n")
			if w.readDefaultString("n") != "n" {
				infos.banned = nil
			}
			// Offer the user to explicitly add/remove certain IP addresses
			fmt.Println()
			fmt.Println("Which additional IP addresses should be blacklisted?")
			for {
				if ip := w.readIPAddress(); ip != "" {
					infos.banned = append(infos.banned, ip)
					continue
				}
				break
			}
			fmt.Println()
			fmt.Println("Which IP addresses should not be blacklisted?")
			for {
				if ip := w.readIPAddress(); ip != "" {
					for i, addr := range infos.banned {
						if ip == addr {
							infos.banned = append(infos.banned[:i], infos.banned[i+1:]...)
							break
						}
					}
					continue
				}
				break
			}
			sort.Strings(infos.banned)
		}
	}
	// Try to deploy the yocstats server on the host
	nocache := false
	if existed {
		fmt.Println()
		fmt.Printf("Should the yocstats be built from scratch (y/n)? (default = no)\n")
		nocache = w.readDefaultString("n") != "n"
	}
	trusted := make([]string, 0, len(w.servers))
	for _, client := range w.servers {
		if client != nil {
			trusted = append(trusted, client.address)
		}
	}
	if out, err := deployEthstats(client, w.network, infos.port, infos.secret, infos.host, trusted, infos.banned, nocache); err != nil {
		log.Error("Failed to deploy yocstats container", "err", err)
		if len(out) > 0 {
			fmt.Printf("%s\n", out)
		}
		return
	}
	// All ok, run a network scan to pick any changes up
	w.networkStats()
}
