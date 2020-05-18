// Authored and revised by YOC team, 2016-2018
// License placeholder #1

// +build !go1.6

package debug

// LoudPanic panics in a way that gets all goroutine stacks printed on stderr.
func LoudPanic(x interface{}) {
	panic(x)
}
