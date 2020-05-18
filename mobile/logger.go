// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package geth

import (
	"os"

	"github.com/Yocoin15/Yocoin_Sources/log"
)

// SetVerbosity sets the global verbosity level (between 0 and 6 - see logger/verbosity.go).
func SetVerbosity(level int) {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(level), log.StreamHandler(os.Stderr, log.TerminalFormat(false))))
}
