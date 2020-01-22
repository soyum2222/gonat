package sign

import (
	"encoding/binary"
	"math"
)

// signature
// sum of all the bytes to remainder 2^31-1 ,and push in the first four byte

func Signature(b []byte) []byte {
	var sum uint32
	signBytes := make([]byte, 4)
	for _, v := range b {
		sum += uint32(v)
	}

	sign := sum % math.MaxUint32

	binary.BigEndian.PutUint32(signBytes[0:4], sign)

	return append(signBytes, b...)

}

func Verifi(b []byte) bool {
	if len(b) < 4 {
		return false
	}

	var sum uint32
	for _, v := range b[4:] {
		sum += uint32(v)
	}

	return sum%math.MaxUint32 == binary.BigEndian.Uint32(b[0:4])
}
