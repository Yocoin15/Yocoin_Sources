// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package yocclient

import (
	"errors"
	"math/big"

	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/core/types"
)

// senderFromServer is a types.Signer that remembers the sender address returned by the RPC
// server. It is stored in the transaction's sender address cache to avoid an additional
// request in TransactionSender.
type senderFromServer struct {
	addr      common.Address
	blockhash common.Hash
}

var errNotCached = errors.New("sender not cached")

func setSenderFromServer(tx *types.Transaction, addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender(&senderFromServer{addr, block}, tx)
}

func (s *senderFromServer) Equal(other types.Signer) bool {
	os, ok := other.(*senderFromServer)
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer) Sender(tx *types.Transaction) (common.Address, error) {
	if s.blockhash == (common.Hash{}) {
		return common.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer) Hash(tx *types.Transaction) common.Hash {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) SignatureValues(tx *types.Transaction, sig []byte) (R, S, V *big.Int, err error) {
	panic("can't sign with senderFromServer")
}
