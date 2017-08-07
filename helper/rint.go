package helper // import "go.rls.moe/secl/helper"

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"encoding/base64"
)

const (
	intByteReadSize = 8
)

// RndFloat uses crypto/rand to obtain a fully random float
func RndFloat() float64 {
	// I really hope this is correct, but it should be mostly fine
	return float64(RndInt64()) / float64(^uint(0))
}

// RndInt64 returns a fully cryptographically random uint64 value
func RndInt64() uint64 {
	var buf = make([]byte, intByteReadSize)
	n, err := rand.Read(buf)
	if err != nil {
		panic(fmt.Sprintf("RndFloat failed, error in read: %s", err.Error()))
	}
	if n < intByteReadSize {
		panic(fmt.Sprintf("RndFloat failed, did not read %d characters, read %d instead", intByteReadSize, n))
	}
	return binary.LittleEndian.Uint64(buf)
}

// RndStr returns a URL safe, cryptographically random string of length n
func RndStr(n int) string {
	if n > 512 {
		// We do this to prevent potential attacks by using excessively sized random strings
		panic("Only 512 RndStr is allowed")
	}
	enc := base64.RawURLEncoding
	plen := enc.DecodedLen(n + 1)
	var buf = make([]byte, plen)
	m, err := rand.Read(buf)
	if err != nil {
		panic(fmt.Sprintf("RndStr failed, error in read: %s", err.Error()))
	}
	if m < plen {
		panic(fmt.Sprintf("RndStr failed, did not read %d characters, read %d instead", plen, m))
	}
	return base64.RawURLEncoding.EncodeToString(buf)[:n]
}
