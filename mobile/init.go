// Authored and revised by YOC team, 2016-2018
// License placeholder #1

// Contains initialization code for the mbile library.

package geth

import (
	"os"
	"runtime"

	"github.com/Yocoin15/Yocoin_Sources/log"
)

func init() {
	// Initialize the logger
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

	// Initialize the goroutine count
	runtime.GOMAXPROCS(runtime.NumCPU())
}
