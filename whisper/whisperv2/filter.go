// Authored and revised by YOC team, 2014-2018
// License placeholder #1

// Contains the message filter for fine grained subscriptions.

package whisperv2

import (
	"crypto/ecdsa"

	"github.com/Yocoin15/Yocoin_Sources/event/filter"
)

// Filter is used to subscribe to specific types of whisper messages.
type Filter struct {
	To     *ecdsa.PublicKey   // Recipient of the message
	From   *ecdsa.PublicKey   // Sender of the message
	Topics [][]Topic          // Topics to filter messages with
	Fn     func(msg *Message) // Handler in case of a match
}

// NewFilterTopics creates a 2D topic array used by whisper.Filter from binary
// data elements.
func NewFilterTopics(data ...[][]byte) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, condition := range data {
		// Handle the special case of condition == [[]byte{}]
		if len(condition) == 1 && len(condition[0]) == 0 {
			filter[i] = []Topic{}
			continue
		}
		// Otherwise flatten normally
		filter[i] = NewTopics(condition...)
	}
	return filter
}

// NewFilterTopicsFlat creates a 2D topic array used by whisper.Filter from flat
// binary data elements.
func NewFilterTopicsFlat(data ...[]byte) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, element := range data {
		// Only add non-wildcard topics
		filter[i] = make([]Topic, 0, 1)
		if len(element) > 0 {
			filter[i] = append(filter[i], NewTopic(element))
		}
	}
	return filter
}

// NewFilterTopicsFromStrings creates a 2D topic array used by whisper.Filter
// from textual data elements.
func NewFilterTopicsFromStrings(data ...[]string) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, condition := range data {
		// Handle the special case of condition == [""]
		if len(condition) == 1 && condition[0] == "" {
			filter[i] = []Topic{}
			continue
		}
		// Otherwise flatten normally
		filter[i] = NewTopicsFromStrings(condition...)
	}
	return filter
}

// NewFilterTopicsFromStringsFlat creates a 2D topic array used by whisper.Filter from flat
// binary data elements.
func NewFilterTopicsFromStringsFlat(data ...string) [][]Topic {
	filter := make([][]Topic, len(data))
	for i, element := range data {
		// Only add non-wildcard topics
		filter[i] = make([]Topic, 0, 1)
		if element != "" {
			filter[i] = append(filter[i], NewTopicFromString(element))
		}
	}
	return filter
}

// filterer is the internal, fully initialized filter ready to match inbound
// messages to a variety of criteria.
type filterer struct {
	to      string                 // Recipient of the message
	from    string                 // Sender of the message
	matcher *topicMatcher          // Topics to filter messages with
	fn      func(data interface{}) // Handler in case of a match
}

// Compare checks if the specified filter matches the current one.
func (self filterer) Compare(f filter.Filter) bool {
	filter := f.(filterer)

	// Check the message sender and recipient
	if len(self.to) > 0 && self.to != filter.to {
		return false
	}
	if len(self.from) > 0 && self.from != filter.from {
		return false
	}
	// Check the topic filtering
	topics := make([]Topic, len(filter.matcher.conditions))
	for i, group := range filter.matcher.conditions {
		// Message should contain a single topic entry, extract
		for topics[i] = range group {
			break
		}
	}
	return self.matcher.Matches(topics)
}

// Trigger is called when a filter successfully matches an inbound message.
func (self filterer) Trigger(data interface{}) {
	self.fn(data)
}
