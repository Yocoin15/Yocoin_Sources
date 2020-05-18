// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main YoCoin network.
var MainnetBootnodes = []string{
	"enode://e246081a1c4418b868f398ef0006ca04d491f0e235cdc23dcc4148cd66e1057500c725b115002876bb83827b7a7aa586406a64fb796c919ffd7e2e1e8030aecb@161.35.106.183:30303",
	"enode://110b389f100045dec04f23f336fd2b099723d7a933a283516a2aad228109dfa8f4793e7a87d205b79271cb724d863ba8ca8a1e2e192d5d9c3ddb2e67934e4e0c@161.35.106.192:30303",
	"enode://fdafd8c9972637b4be79548b6be1379637147b314ea8911236d7901a73596dfcb37f64864db4b56d1fbadcff6199d60ad5f2fa4173fe3e57b41c753946c56ada@68.183.92.142:30303",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://e246081a1c4418b868f398ef0006ca04d491f0e235cdc23dcc4148cd66e1057500c725b115002876bb83827b7a7aa586406a64fb796c919ffd7e2e1e8030aecb@161.35.106.183:30303",
	"enode://110b389f100045dec04f23f336fd2b099723d7a933a283516a2aad228109dfa8f4793e7a87d205b79271cb724d863ba8ca8a1e2e192d5d9c3ddb2e67934e4e0c@161.35.106.192:30303",
	"enode://fdafd8c9972637b4be79548b6be1379637147b314ea8911236d7901a73596dfcb37f64864db4b56d1fbadcff6199d60ad5f2fa4173fe3e57b41c753946c56ada@68.183.92.142:30303",
}
