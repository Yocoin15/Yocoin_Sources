// Authored and revised by YOC team, 2016-2018
// License placeholder #1

// package mclock is a wrapper for a monotonic clock source
package mclock

import (
	"time"

	"github.com/aristanetworks/goarista/monotime"
)

type AbsTime time.Duration // absolute monotonic time

func Now() AbsTime {
	return AbsTime(monotime.Now())
}
