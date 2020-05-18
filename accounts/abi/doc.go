// Authored and revised by YOC team, 2015-2018
// License placeholder #1

// Package abi implements the YOC ABI (Application Binary
// Interface).
//
// The YOC ABI is strongly typed, known at compile time
// and static. This ABI will handle basic type casting; unsigned
// to signed and visa versa. It does not handle slice casting such
// as unsigned slice to signed slice. Bit size type casting is also
// handled. ints with a bit size of 32 will be properly cast to int256,
// etc.
package abi
