// Authored and revised by YOC team, 2016-2018
// License placeholder #1

//+build !windows

package netutil

// isPacketTooBig reports whether err indicates that a UDP packet didn't
// fit the receive buffer. There is no such error on
// non-Windows platforms.
func isPacketTooBig(err error) bool {
	return false
}
