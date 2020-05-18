// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package api

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/contracts/ens"
	"github.com/Yocoin15/Yocoin_Sources/crypto"
	"github.com/Yocoin15/Yocoin_Sources/log"
	"github.com/Yocoin15/Yocoin_Sources/node"
	"github.com/Yocoin15/Yocoin_Sources/swarm/network"
	"github.com/Yocoin15/Yocoin_Sources/swarm/services/swap"
	"github.com/Yocoin15/Yocoin_Sources/swarm/storage"
)

const (
	DefaultHTTPListenAddr = "127.0.0.1"
	DefaultHTTPPort       = "8500"
)

// separate bzz directories
// allow several bzz nodes running in parallel
type Config struct {
	// serialised/persisted fields
	*storage.StoreParams
	*storage.ChunkerParams
	*network.HiveParams
	Swap        *swap.SwapParams
	*network.SyncParams
	Contract    common.Address
	EnsRoot     common.Address
	EnsApi      string
	Path        string
	ListenAddr  string
	Port        string
	PublicKey   string
	BzzKey      string
	NetworkId   uint64
	SwapEnabled bool
	SyncEnabled bool
	SwapApi     string
	Cors        string
	BzzAccount  string
	BootNodes   string
}

//create a default config with all parameters to set to defaults
func NewDefaultConfig() (self *Config) {

	self = &Config{
		StoreParams:   storage.NewDefaultStoreParams(),
		ChunkerParams: storage.NewChunkerParams(),
		HiveParams:    network.NewDefaultHiveParams(),
		SyncParams:    network.NewDefaultSyncParams(),
		Swap:          swap.NewDefaultSwapParams(),
		ListenAddr:    DefaultHTTPListenAddr,
		Port:          DefaultHTTPPort,
		Path:          node.DefaultDataDir(),
		EnsApi:        node.DefaultIPCEndpoint("geth"),
		EnsRoot:       ens.TestNetAddress,
		NetworkId:     network.NetworkId,
		SwapEnabled:   false,
		SyncEnabled:   true,
		SwapApi:       "",
		BootNodes:     "",
	}

	return
}

//some config params need to be initialized after the complete
//config building phase is completed (e.g. due to overriding flags)
func (self *Config) Init(prvKey *ecdsa.PrivateKey) {

	address := crypto.PubkeyToAddress(prvKey.PublicKey)
	self.Path = filepath.Join(self.Path, "bzz-"+common.Bytes2Hex(address.Bytes()))
	err := os.MkdirAll(self.Path, os.ModePerm)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating root swarm data directory: %v", err))
		return
	}

	pubkey := crypto.FromECDSAPub(&prvKey.PublicKey)
	pubkeyhex := common.ToHex(pubkey)
	keyhex := crypto.Keccak256Hash(pubkey).Hex()

	self.PublicKey = pubkeyhex
	self.BzzKey = keyhex

	self.Swap.Init(self.Contract, prvKey)
	self.SyncParams.Init(self.Path)
	self.HiveParams.Init(self.Path)
	self.StoreParams.Init(self.Path)
}
