// Authored and revised by YOC team, 2014-2018
// License placeholder #1

// Package rle implements the run-length encoding used for YOC data.
package rle

import (
	"bytes"
	"errors"

	"github.com/Yocoin15/Yocoin_Sources/crypto"
)

const (
	token             byte = 0xfe
	emptyShaToken          = 0xfd
	emptyListShaToken      = 0xfe
	tokenToken             = 0xff
)

var empty = crypto.Keccak256([]byte(""))
var emptyList = crypto.Keccak256([]byte{0x80})

func Decompress(dat []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	for i := 0; i < len(dat); i++ {
		if dat[i] == token {
			if i+1 < len(dat) {
				switch dat[i+1] {
				case emptyShaToken:
					buf.Write(empty)
				case emptyListShaToken:
					buf.Write(emptyList)
				case tokenToken:
					buf.WriteByte(token)
				default:
					buf.Write(make([]byte, int(dat[i+1]-2)))
				}
				i++
			} else {
				return nil, errors.New("error reading bytes. token encountered without proceeding bytes")
			}
		} else {
			buf.WriteByte(dat[i])
		}
	}

	return buf.Bytes(), nil
}

func compressChunk(dat []byte) (ret []byte, n int) {
	switch {
	case dat[0] == token:
		return []byte{token, tokenToken}, 1
	case len(dat) > 1 && dat[0] == 0x0 && dat[1] == 0x0:
		j := 0
		for j <= 254 && j < len(dat) {
			if dat[j] != 0 {
				break
			}
			j++
		}
		return []byte{token, byte(j + 2)}, j
	case len(dat) >= 32:
		if dat[0] == empty[0] && bytes.Equal(dat[:32], empty) {
			return []byte{token, emptyShaToken}, 32
		} else if dat[0] == emptyList[0] && bytes.Equal(dat[:32], emptyList) {
			return []byte{token, emptyListShaToken}, 32
		}
		fallthrough
	default:
		return dat[:1], 1
	}
}

func Compress(dat []byte) []byte {
	buf := new(bytes.Buffer)

	i := 0
	for i < len(dat) {
		b, n := compressChunk(dat[i:])
		buf.Write(b)
		i += n
	}

	return buf.Bytes()
}
