// Authored and revised by YOC team, 2017-2018
// License placeholder #1

// +build linux darwin freebsd

package fuse

import (
	"bazil.org/fuse/fs"
)

var (
	_ fs.Node = (*SwarmDir)(nil)
)

type SwarmRoot struct {
	root *SwarmDir
}

func (filesystem *SwarmRoot) Root() (fs.Node, error) {
	return filesystem.root, nil
}
