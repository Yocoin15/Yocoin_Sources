// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package abi

import (
	"fmt"
	"strings"

	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/crypto"
)

// Event is an event potentially triggered by the YVM's LOG mechanism. The Event
// holds type information (inputs) about the yielded output. Anonymous events
// don't get the signature canonical representation as the first LOG topic.
type Event struct {
	Name      string
	Anonymous bool
	Inputs    Arguments
}

func (event Event) String() string {
	inputs := make([]string, len(event.Inputs))
	for i, input := range event.Inputs {
		inputs[i] = fmt.Sprintf("%v %v", input.Name, input.Type)
		if input.Indexed {
			inputs[i] = fmt.Sprintf("%v indexed %v", input.Name, input.Type)
		}
	}
	return fmt.Sprintf("event %v(%v)", event.Name, strings.Join(inputs, ", "))
}

// Id returns the canonical representation of the event's signature used by the
// abi definition to identify event names and types.
func (e Event) Id() common.Hash {
	types := make([]string, len(e.Inputs))
	i := 0
	for _, input := range e.Inputs {
		types[i] = input.Type.String()
		i++
	}
	return common.BytesToHash(crypto.Keccak256([]byte(fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ",")))))
}
