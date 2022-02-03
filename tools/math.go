package tools

import (
	"encoding/binary"
	mrand "math/rand"
	crand "crypto/rand"
)

func RandomSeed() {
	var b [8]byte

	_, err := crand.Read(b[:])

	if err != nil {
		panic("failed to seed math/rand")
	}

	mrand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
