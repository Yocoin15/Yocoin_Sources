// Authored and revised by YOC team, 2018
// License placeholder #1

package dashboard

// getProcessCPUTime returns 0 on Windows as there is no system call to resolve
// the actual process' CPU time.
func getProcessCPUTime() float64 {
	return 0
}
