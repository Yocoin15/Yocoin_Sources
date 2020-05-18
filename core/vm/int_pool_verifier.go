// Authored and revised by YOC team, 2017-2018
// License placeholder #1

// +build VERIFY_YVM_INTEGER_POOL

package vm

import "fmt"

const verifyPool = true

func verifyIntegerPool(ip *intPool) {
	for i, item := range ip.pool.data {
		if item.Cmp(checkVal) != 0 {
			panic(fmt.Sprintf("%d'th item failed aggressive pool check. Value was modified", i))
		}
	}
}
