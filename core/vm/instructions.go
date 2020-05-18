// Authored and revised by YOC team, 2015-2018
// License placeholder #1

package vm

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/common/math"
	"github.com/Yocoin15/Yocoin_Sources/core/types"
	"github.com/Yocoin15/Yocoin_Sources/crypto"
	"github.com/Yocoin15/Yocoin_Sources/params"
)

var (
	bigZero                  = new(big.Int)
	errWriteProtection       = errors.New("evm: write protection")
	errReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	errExecutionReverted     = errors.New("evm: execution reverted")
	errMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")
)

func opAdd(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Add(x, y)))

	yvm.interpreter.intPool.put(y)

	return nil, nil
}

func opSub(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Sub(x, y)))

	yvm.interpreter.intPool.put(y)

	return nil, nil
}

func opMul(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(math.U256(x.Mul(x, y)))

	yvm.interpreter.intPool.put(y)

	return nil, nil
}

func opDiv(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	if y.Sign() != 0 {
		stack.push(math.U256(x.Div(x, y)))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(y)

	return nil, nil
}

func opSdiv(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if y.Sign() == 0 {
		stack.push(new(big.Int))
		return nil, nil
	} else {
		n := new(big.Int)
		if yvm.interpreter.intPool.get().Mul(x, y).Sign() < 0 {
			n.SetInt64(-1)
		} else {
			n.SetInt64(1)
		}

		res := x.Div(x.Abs(x), y.Abs(y))
		res.Mul(res, n)

		stack.push(math.U256(res))
	}
	yvm.interpreter.intPool.put(y)
	return nil, nil
}

func opMod(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	if y.Sign() == 0 {
		stack.push(new(big.Int))
	} else {
		stack.push(math.U256(x.Mod(x, y)))
	}
	yvm.interpreter.intPool.put(y)
	return nil, nil
}

func opSmod(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := math.S256(stack.pop()), math.S256(stack.pop())

	if y.Sign() == 0 {
		stack.push(new(big.Int))
	} else {
		n := new(big.Int)
		if x.Sign() < 0 {
			n.SetInt64(-1)
		} else {
			n.SetInt64(1)
		}

		res := x.Mod(x.Abs(x), y.Abs(y))
		res.Mul(res, n)

		stack.push(math.U256(res))
	}
	yvm.interpreter.intPool.put(y)
	return nil, nil
}

func opExp(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	base, exponent := stack.pop(), stack.pop()
	stack.push(math.Exp(base, exponent))

	yvm.interpreter.intPool.put(base, exponent)

	return nil, nil
}

func opSignExtend(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	back := stack.pop()
	if back.Cmp(big.NewInt(31)) < 0 {
		bit := uint(back.Uint64()*8 + 7)
		num := stack.pop()
		mask := back.Lsh(common.Big1, bit)
		mask.Sub(mask, common.Big1)
		if num.Bit(int(bit)) > 0 {
			num.Or(num, mask.Not(mask))
		} else {
			num.And(num, mask)
		}

		stack.push(math.U256(num))
	}

	yvm.interpreter.intPool.put(back)
	return nil, nil
}

func opNot(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x := stack.pop()
	stack.push(math.U256(x.Not(x)))
	return nil, nil
}

func opLt(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) < 0 {
		stack.push(yvm.interpreter.intPool.get().SetUint64(1))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opGt(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) > 0 {
		stack.push(yvm.interpreter.intPool.get().SetUint64(1))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opSlt(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if x.Cmp(math.S256(y)) < 0 {
		stack.push(yvm.interpreter.intPool.get().SetUint64(1))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opSgt(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := math.S256(stack.pop()), math.S256(stack.pop())
	if x.Cmp(y) > 0 {
		stack.push(yvm.interpreter.intPool.get().SetUint64(1))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opEq(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	if x.Cmp(y) == 0 {
		stack.push(yvm.interpreter.intPool.get().SetUint64(1))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(x, y)
	return nil, nil
}

func opIszero(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x := stack.pop()
	if x.Sign() > 0 {
		stack.push(new(big.Int))
	} else {
		stack.push(yvm.interpreter.intPool.get().SetUint64(1))
	}

	yvm.interpreter.intPool.put(x)
	return nil, nil
}

func opAnd(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(x.And(x, y))

	yvm.interpreter.intPool.put(y)
	return nil, nil
}

func opOr(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(x.Or(x, y))

	yvm.interpreter.intPool.put(y)
	return nil, nil
}

func opXor(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y := stack.pop(), stack.pop()
	stack.push(x.Xor(x, y))

	yvm.interpreter.intPool.put(y)
	return nil, nil
}

func opByte(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	th, val := stack.pop(), stack.peek()
	if th.Cmp(common.Big32) < 0 {
		b := math.Byte(val, 32, int(th.Int64()))
		val.SetUint64(uint64(b))
	} else {
		val.SetUint64(0)
	}
	yvm.interpreter.intPool.put(th)
	return nil, nil
}

func opAddmod(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		add := x.Add(x, y)
		add.Mod(add, z)
		stack.push(math.U256(add))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opMulmod(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	x, y, z := stack.pop(), stack.pop(), stack.pop()
	if z.Cmp(bigZero) > 0 {
		mul := x.Mul(x, y)
		mul.Mod(mul, z)
		stack.push(math.U256(mul))
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(y, z)
	return nil, nil
}

func opSha3(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	data := memory.Get(offset.Int64(), size.Int64())
	hash := crypto.Keccak256(data)

	if yvm.vmConfig.EnablePreimageRecording {
		yvm.StateDB.AddPreimage(common.BytesToHash(hash), data)
	}

	stack.push(new(big.Int).SetBytes(hash))

	yvm.interpreter.intPool.put(offset, size)
	return nil, nil
}

func opAddress(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(contract.Address().Big())
	return nil, nil
}

func opBalance(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	addr := common.BigToAddress(stack.pop())
	balance := yvm.StateDB.GetBalance(addr)

	stack.push(new(big.Int).Set(balance))
	return nil, nil
}

func opOrigin(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.Origin.Big())
	return nil, nil
}

func opCaller(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(contract.Caller().Big())
	return nil, nil
}

func opCallValue(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().Set(contract.value))
	return nil, nil
}

func opCallDataLoad(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(new(big.Int).SetBytes(getDataBig(contract.Input, stack.pop(), big32)))
	return nil, nil
}

func opCallDataSize(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().SetInt64(int64(len(contract.Input))))
	return nil, nil
}

func opCallDataCopy(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	memory.Set(memOffset.Uint64(), length.Uint64(), getDataBig(contract.Input, dataOffset, length))

	yvm.interpreter.intPool.put(memOffset, dataOffset, length)
	return nil, nil
}

func opReturnDataSize(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().SetUint64(uint64(len(yvm.interpreter.returnData))))
	return nil, nil
}

func opReturnDataCopy(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		dataOffset = stack.pop()
		length     = stack.pop()
	)
	defer yvm.interpreter.intPool.put(memOffset, dataOffset, length)

	end := new(big.Int).Add(dataOffset, length)
	if end.BitLen() > 64 || uint64(len(yvm.interpreter.returnData)) < end.Uint64() {
		return nil, errReturnDataOutOfBounds
	}
	memory.Set(memOffset.Uint64(), length.Uint64(), yvm.interpreter.returnData[dataOffset.Uint64():end.Uint64()])

	return nil, nil
}

func opExtCodeSize(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	a := stack.pop()

	addr := common.BigToAddress(a)
	a.SetInt64(int64(yvm.StateDB.GetCodeSize(addr)))
	stack.push(a)

	return nil, nil
}

func opCodeSize(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	l := yvm.interpreter.intPool.get().SetInt64(int64(len(contract.Code)))
	stack.push(l)
	return nil, nil
}

func opCodeCopy(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := getDataBig(contract.Code, codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	yvm.interpreter.intPool.put(memOffset, codeOffset, length)
	return nil, nil
}

func opExtCodeCopy(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		addr       = common.BigToAddress(stack.pop())
		memOffset  = stack.pop()
		codeOffset = stack.pop()
		length     = stack.pop()
	)
	codeCopy := getDataBig(yvm.StateDB.GetCode(addr), codeOffset, length)
	memory.Set(memOffset.Uint64(), length.Uint64(), codeCopy)

	yvm.interpreter.intPool.put(memOffset, codeOffset, length)
	return nil, nil
}

func opGasprice(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().Set(yvm.GasPrice))
	return nil, nil
}

func opBlockhash(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	num := stack.pop()

	n := yvm.interpreter.intPool.get().Sub(yvm.BlockNumber, common.Big257)
	if num.Cmp(n) > 0 && num.Cmp(yvm.BlockNumber) < 0 {
		stack.push(yvm.GetHash(num.Uint64()).Big())
	} else {
		stack.push(new(big.Int))
	}

	yvm.interpreter.intPool.put(num, n)
	return nil, nil
}

func opCoinbase(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.Coinbase.Big())
	return nil, nil
}

func opTimestamp(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(math.U256(new(big.Int).Set(yvm.Time)))
	return nil, nil
}

func opNumber(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(math.U256(new(big.Int).Set(yvm.BlockNumber)))
	return nil, nil
}

func opDifficulty(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(math.U256(new(big.Int).Set(yvm.Difficulty)))
	return nil, nil
}

func opGasLimit(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(math.U256(new(big.Int).SetUint64(yvm.GasLimit)))
	return nil, nil
}

func opPop(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	yvm.interpreter.intPool.put(stack.pop())
	return nil, nil
}

func opMload(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset := stack.pop()
	val := new(big.Int).SetBytes(memory.Get(offset.Int64(), 32))
	stack.push(val)

	yvm.interpreter.intPool.put(offset)
	return nil, nil
}

func opMstore(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// pop value of the stack
	mStart, val := stack.pop(), stack.pop()
	memory.Set(mStart.Uint64(), 32, math.PaddedBigBytes(val, 32))

	yvm.interpreter.intPool.put(mStart, val)
	return nil, nil
}

func opMstore8(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	off, val := stack.pop().Int64(), stack.pop().Int64()
	memory.store[off] = byte(val & 0xff)

	return nil, nil
}

func opSload(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	loc := common.BigToHash(stack.pop())
	val := yvm.StateDB.GetState(contract.Address(), loc).Big()
	stack.push(val)
	return nil, nil
}

func opSstore(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	loc := common.BigToHash(stack.pop())
	val := stack.pop()
	yvm.StateDB.SetState(contract.Address(), loc, common.BigToHash(val))

	yvm.interpreter.intPool.put(val)
	return nil, nil
}

func opJump(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	pos := stack.pop()
	if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos) {
		nop := contract.GetOp(pos.Uint64())
		return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
	}
	*pc = pos.Uint64()

	yvm.interpreter.intPool.put(pos)
	return nil, nil
}

func opJumpi(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	pos, cond := stack.pop(), stack.pop()
	if cond.Sign() != 0 {
		if !contract.jumpdests.has(contract.CodeHash, contract.Code, pos) {
			nop := contract.GetOp(pos.Uint64())
			return nil, fmt.Errorf("invalid jump destination (%v) %v", nop, pos)
		}
		*pc = pos.Uint64()
	} else {
		*pc++
	}

	yvm.interpreter.intPool.put(pos, cond)
	return nil, nil
}

func opJumpdest(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	return nil, nil
}

func opPc(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().SetUint64(*pc))
	return nil, nil
}

func opMsize(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().SetInt64(int64(memory.Len())))
	return nil, nil
}

func opGas(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	stack.push(yvm.interpreter.intPool.get().SetUint64(contract.Gas))
	return nil, nil
}

func opCreate(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	var (
		value        = stack.pop()
		offset, size = stack.pop(), stack.pop()
		input        = memory.Get(offset.Int64(), size.Int64())
		gas          = contract.Gas
	)
	if yvm.ChainConfig().IsEIP150(yvm.BlockNumber) {
		gas -= gas / 64
	}

	contract.UseGas(gas)
	res, addr, returnGas, suberr := yvm.Create(contract, input, gas, value)
	// Push item on the stack based on the returned error. If the ruleset is
	// homestead we must check for CodeStoreOutOfGasError (homestead only
	// rule) and treat as an error, if the ruleset is frontier we must
	// ignore this error and pretend the operation was successful.
	if yvm.ChainConfig().IsHomestead(yvm.BlockNumber) && suberr == ErrCodeStoreOutOfGas {
		stack.push(new(big.Int))
	} else if suberr != nil && suberr != ErrCodeStoreOutOfGas {
		stack.push(new(big.Int))
	} else {
		stack.push(addr.Big())
	}
	contract.Gas += returnGas
	yvm.interpreter.intPool.put(value, offset, size)

	if suberr == errExecutionReverted {
		return res, nil
	}
	return nil, nil
}

func opCall(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas in in yvm.callGasTemp.
	yvm.interpreter.intPool.put(stack.pop())
	gas := yvm.callGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.BigToAddress(addr)
	value = math.U256(value)
	// Get the arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := yvm.Call(contract, toAddr, args, gas, value)
	if err != nil {
		stack.push(new(big.Int))
	} else {
		stack.push(big.NewInt(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	yvm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opCallCode(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in yvm.callGasTemp.
	yvm.interpreter.intPool.put(stack.pop())
	gas := yvm.callGasTemp
	// Pop other call parameters.
	addr, value, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.BigToAddress(addr)
	value = math.U256(value)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	ret, returnGas, err := yvm.CallCode(contract, toAddr, args, gas, value)
	if err != nil {
		stack.push(new(big.Int))
	} else {
		stack.push(big.NewInt(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	yvm.interpreter.intPool.put(addr, value, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opDelegateCall(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in yvm.callGasTemp.
	yvm.interpreter.intPool.put(stack.pop())
	gas := yvm.callGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.BigToAddress(addr)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := yvm.DelegateCall(contract, toAddr, args, gas)
	if err != nil {
		stack.push(new(big.Int))
	} else {
		stack.push(big.NewInt(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	yvm.interpreter.intPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opStaticCall(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	// Pop gas. The actual gas is in yvm.callGasTemp.
	yvm.interpreter.intPool.put(stack.pop())
	gas := yvm.callGasTemp
	// Pop other call parameters.
	addr, inOffset, inSize, retOffset, retSize := stack.pop(), stack.pop(), stack.pop(), stack.pop(), stack.pop()
	toAddr := common.BigToAddress(addr)
	// Get arguments from the memory.
	args := memory.Get(inOffset.Int64(), inSize.Int64())

	ret, returnGas, err := yvm.StaticCall(contract, toAddr, args, gas)
	if err != nil {
		stack.push(new(big.Int))
	} else {
		stack.push(big.NewInt(1))
	}
	if err == nil || err == errExecutionReverted {
		memory.Set(retOffset.Uint64(), retSize.Uint64(), ret)
	}
	contract.Gas += returnGas

	yvm.interpreter.intPool.put(addr, inOffset, inSize, retOffset, retSize)
	return ret, nil
}

func opReturn(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	yvm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opRevert(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	offset, size := stack.pop(), stack.pop()
	ret := memory.GetPtr(offset.Int64(), size.Int64())

	yvm.interpreter.intPool.put(offset, size)
	return ret, nil
}

func opStop(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	return nil, nil
}

func opSuicide(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
	balance := yvm.StateDB.GetBalance(contract.Address())
	yvm.StateDB.AddBalance(common.BigToAddress(stack.pop()), balance)

	yvm.StateDB.Suicide(contract.Address())
	return nil, nil
}

// following functions are used by the instruction jump  table

// make log instruction function
func makeLog(size int) executionFunc {
	return func(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		topics := make([]common.Hash, size)
		mStart, mSize := stack.pop(), stack.pop()
		for i := 0; i < size; i++ {
			topics[i] = common.BigToHash(stack.pop())
		}

		d := memory.Get(mStart.Int64(), mSize.Int64())
		yvm.StateDB.AddLog(&types.Log{
			Address: contract.Address(),
			Topics:  topics,
			Data:    d,
			// This is a non-consensus field, but assigned here because
			// core/state doesn't know the current block number.
			BlockNumber: yvm.BlockNumber.Uint64(),
		})

		yvm.interpreter.intPool.put(mStart, mSize)
		return nil, nil
	}
}

// make push instruction function
func makePush(size uint64, pushByteSize int) executionFunc {
	return func(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		codeLen := len(contract.Code)

		startMin := codeLen
		if int(*pc+1) < startMin {
			startMin = int(*pc + 1)
		}

		endMin := codeLen
		if startMin+pushByteSize < endMin {
			endMin = startMin + pushByteSize
		}

		integer := yvm.interpreter.intPool.get()
		stack.push(integer.SetBytes(common.RightPadBytes(contract.Code[startMin:endMin], pushByteSize)))

		*pc += size
		return nil, nil
	}
}

// make push instruction function
func makeDup(size int64) executionFunc {
	return func(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		stack.dup(yvm.interpreter.intPool, int(size))
		return nil, nil
	}
}

// make swap instruction function
func makeSwap(size int64) executionFunc {
	// switch n + 1 otherwise n would be swapped with n
	size += 1
	return func(pc *uint64, yvm *YVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error) {
		stack.swap(int(size))
		return nil, nil
	}
}
