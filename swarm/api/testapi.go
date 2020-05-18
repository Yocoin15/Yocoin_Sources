// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package api

import (
	"github.com/Yocoin15/Yocoin_Sources/swarm/network"
)

type Control struct {
	api  *Api
	hive *network.Hive
}

func NewControl(api *Api, hive *network.Hive) *Control {
	return &Control{api, hive}
}

func (self *Control) BlockNetworkRead(on bool) {
	self.hive.BlockNetworkRead(on)
}

func (self *Control) SyncEnabled(on bool) {
	self.hive.SyncEnabled(on)
}

func (self *Control) SwapEnabled(on bool) {
	self.hive.SwapEnabled(on)
}

func (self *Control) Hive() string {
	return self.hive.String()
}
