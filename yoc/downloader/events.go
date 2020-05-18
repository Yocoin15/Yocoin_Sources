// Authored and revised by YOC team, 2015-2018
// License placeholder #1

package downloader

type DoneEvent struct{}
type StartEvent struct{}
type FailedEvent struct{ Err error }
