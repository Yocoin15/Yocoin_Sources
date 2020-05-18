// Authored and revised by YOC team, 2016-2018
// License placeholder #1

// Contains the Whisper protocol Topic element.

package whisperv6

import (
	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/common/hexutil"
)

// Topic represents a cryptographically secure, probabilistic partial
// classifications of a message, determined as the first (left) 4 bytes of the
// SHA3 hash of some arbitrary data given by the original author of the message.
type TopicType [TopicLength]byte

func BytesToTopic(b []byte) (t TopicType) {
	sz := TopicLength
	if x := len(b); x < TopicLength {
		sz = x
	}
	for i := 0; i < sz; i++ {
		t[i] = b[i]
	}
	return t
}

// String converts a topic byte array to a string representation.
func (t *TopicType) String() string {
	return common.ToHex(t[:])
}

// MarshalText returns the hex representation of t.
func (t TopicType) MarshalText() ([]byte, error) {
	return hexutil.Bytes(t[:]).MarshalText()
}

// UnmarshalText parses a hex representation to a topic.
func (t *TopicType) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Topic", input, t[:])
}
