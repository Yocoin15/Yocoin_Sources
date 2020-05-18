// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package les

import (
	"context"
	"time"

	"github.com/Yocoin15/Yocoin_Sources/core"
	"github.com/Yocoin15/Yocoin_Sources/yoc/downloader"
	"github.com/Yocoin15/Yocoin_Sources/light"
)

const (
	//forceSyncCycle      = 10 * time.Second // Time interval to force syncs, even if few peers are available
	minDesiredPeerCount = 5 // Amount of peers desired to start syncing
)

// syncer is responsible for periodically synchronising with the network, both
// downloading hashes and blocks as well as handling the announcement handler.
func (pm *ProtocolManager) syncer() {
	// Start and ensure cleanup of sync mechanisms
	//pm.fetcher.Start()
	//defer pm.fetcher.Stop()
	defer pm.downloader.Terminate()

	// Wait for different events to fire synchronisation operations
	//forceSync := time.Tick(forceSyncCycle)
	for {
		select {
		case <-pm.newPeerCh:
			/*			// Make sure we have peers to select from, then sync
						if pm.peers.Len() < minDesiredPeerCount {
							break
						}
						go pm.synchronise(pm.peers.BestPeer())
			*/
		/*case <-forceSync:
		// Force a sync even if not enough peers are present
		go pm.synchronise(pm.peers.BestPeer())
		*/
		case <-pm.noMorePeers:
			return
		}
	}
}

func (pm *ProtocolManager) needToSync(peerHead blockInfo) bool {
	head := pm.blockchain.CurrentHeader()
	currentTd := core.GetTd(pm.chainDb, head.Hash(), head.Number.Uint64())
	return currentTd != nil && peerHead.Td.Cmp(currentTd) > 0
}

// synchronise tries to sync up our local block chain with a remote peer.
func (pm *ProtocolManager) synchronise(peer *peer) {
	// Short circuit if no peers are available
	if peer == nil {
		return
	}

	// Make sure the peer's TD is higher than our own.
	if !pm.needToSync(peer.headBlockInfo()) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pm.blockchain.(*light.LightChain).SyncCht(ctx)
	pm.downloader.Synchronise(peer.id, peer.Head(), peer.Td(), downloader.LightSync)
}
