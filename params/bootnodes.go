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
"enode://b0faed40ba59c05b0340798de5e18e2e3dda5e375dd2466fe7ea709c21436376c74ec282513e6ec8a10c79f6cc47228e2b6568732a850e7878b22319cab68cd7@18.159.199.169:30303",
"enode://810ab793ac3d6aae1502fd5727c15a03efe3f8f877c6495cd6a21667fd5c23c4ad3c77e9fa6303e53e855a7daab2b3f91ef94068d85efa12d3d9ddd8df01a50f@18.193.34.193:30303",
"enode://8e5b629d404668d4a6f1a813bd8786cb40e078f1db2e80cfc256a7e9b7ef8430a1a98f78f9afbfde15ba2f00ee5aa04ed0264bd296b78be208c17371b525b48a@52.53.194.227:30303",
"enode://3f38479f09763f0602edb4107c8956d1956c8d857034dca6871f7368351ed95ec3091694f5b440aef0c00093b6219c1daab946ed8860e2351f6944fb3a38fc3b@13.124.37.39:30303",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
"enode://b0faed40ba59c05b0340798de5e18e2e3dda5e375dd2466fe7ea709c21436376c74ec282513e6ec8a10c79f6cc47228e2b6568732a850e7878b22319cab68cd7@18.159.199.169:30303",
"enode://810ab793ac3d6aae1502fd5727c15a03efe3f8f877c6495cd6a21667fd5c23c4ad3c77e9fa6303e53e855a7daab2b3f91ef94068d85efa12d3d9ddd8df01a50f@18.193.34.193:30303",
"enode://8e5b629d404668d4a6f1a813bd8786cb40e078f1db2e80cfc256a7e9b7ef8430a1a98f78f9afbfde15ba2f00ee5aa04ed0264bd296b78be208c17371b525b48a@52.53.194.227:30303",
"enode://3f38479f09763f0602edb4107c8956d1956c8d857034dca6871f7368351ed95ec3091694f5b440aef0c00093b6219c1daab946ed8860e2351f6944fb3a38fc3b@13.124.37.39:30303",
}
