// Authored and revised by YOC team, 2016-2018
// License placeholder #1

package params

import (
	"fmt"
)

const (
	VersionMajor = 1          // Major version component of the current release
	VersionMinor = 8          // Minor version component of the current release
	VersionPatch = 0          // Patch version component of the current release
	VersionMeta  = "unstable" // Version metadata to append to the version string
)

// Version holds the textual version string.
var Version = func() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

func VersionWithCommit(gitCommit string) string {
	vsn := Version
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}
