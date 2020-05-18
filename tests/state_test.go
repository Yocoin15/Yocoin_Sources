// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package tests

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/Yocoin15/Yocoin_Sources/core/vm"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.skipShortMode(`^stQuadraticComplexityTest/`)
	// Broken tests:
	st.skipLoad(`^stTransactionTest/OverflowGasRequire\.json`) // gasLimit > 256 bits
	st.skipLoad(`^stTransactionTest/zeroSigTransa[^/]*\.json`) // EIP-86 is not supported yet
	// Expected failures:
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/EIP158`, "bug in test")
	st.fails(`^stRevertTest/RevertPrefoundEmptyOOG\.json/EIP158`, "bug in test")
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/Byzantium`, "bug in test")
	st.fails(`^stRevertTest/RevertPrefoundEmptyOOG\.json/Byzantium`, "bug in test")
	st.fails(`^stRandom2/randomStatetest64[45]\.json/(EIP150|Frontier|Homestead)/.*`, "known bug #15119")
	st.fails(`^stCreateTest/TransactionCollisionToEmpty\.json/EIP158/2`, "known bug ")
	st.fails(`^stCreateTest/TransactionCollisionToEmpty\.json/EIP158/3`, "known bug ")
	st.fails(`^stCreateTest/TransactionCollisionToEmpty\.json/Byzantium/2`, "known bug ")
	st.fails(`^stCreateTest/TransactionCollisionToEmpty\.json/Byzantium/3`, "known bug ")

	st.walk(t, stateTestDir, func(t *testing.T, name string, test *StateTest) {
		for _, subtest := range test.Subtests() {
			subtest := subtest
			key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
			name := name + "/" + key
			t.Run(key, func(t *testing.T) {
				if subtest.Fork == "Constantinople" {
					t.Skip("constantinople not supported yet")
				}
				withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
					_, err := test.Run(subtest, vmconfig)
					return st.checkFailure(t, name, err)
				})
			})
		}
	})
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	err := test(vm.Config{})
	if err == nil {
		return
	}
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for YVM trace")
		return
	}
	tracer := vm.NewStructLogger(nil)
	err2 := test(vm.Config{Debug: true, Tracer: tracer})
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	buf := new(bytes.Buffer)
	vm.WriteTrace(buf, tracer.StructLogs())
	if buf.Len() == 0 {
		t.Log("no YVM operation logs generated")
	} else {
		t.Log("YVM operation log:\n" + buf.String())
	}
	t.Logf("YVM output: 0x%x", tracer.Output())
	t.Logf("YVM error: %v", tracer.Error())
}
