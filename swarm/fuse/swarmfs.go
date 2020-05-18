// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package fuse

import (
	"github.com/Yocoin15/Yocoin_Sources/swarm/api"
	"sync"
	"time"
)

const (
	Swarmfs_Version = "0.1"
	mountTimeout    = time.Second * 5
	unmountTimeout  = time.Second * 10
	maxFuseMounts   = 5
)

var (
	swarmfs     *SwarmFS // Swarm file system singleton
	swarmfsLock sync.Once

	inode     uint64 = 1 // global inode
	inodeLock sync.RWMutex
)

type SwarmFS struct {
	swarmApi     *api.Api
	activeMounts map[string]*MountInfo
	swarmFsLock  *sync.RWMutex
}

func NewSwarmFS(api *api.Api) *SwarmFS {
	swarmfsLock.Do(func() {
		swarmfs = &SwarmFS{
			swarmApi:     api,
			swarmFsLock:  &sync.RWMutex{},
			activeMounts: map[string]*MountInfo{},
		}
	})
	return swarmfs

}

// Inode numbers need to be unique, they are used for caching inside fuse
func NewInode() uint64 {
	inodeLock.Lock()
	defer inodeLock.Unlock()
	inode += 1
	return inode
}
