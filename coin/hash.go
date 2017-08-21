package coin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/ripemd160"
)

const (
	COIN                     = 1e8
	MAX_MONEY                = 21000000 * COIN
	MAX_BLOCK_SIZE           = 1e6
	MessageMagic             = "Bitcoin Signed Message:\n"
	LOCKTIME_THRESHOLD       = 500000000
	MAX_SCRIPT_ELEMENT_SIZE  = 520
	MAX_BLOCK_SIGOPS_COST    = 80000
	MAX_PUBKEYS_PER_MULTISIG = 20
	WITNESS_SCALE_FACTOR     = 4
)

func ShaHash(b []byte, out []byte) {
	s := sha256.New()
	s.Write(b[:])
	tmp := s.Sum(nil)
	s.Reset()
	s.Write(tmp)
	copy(out[:], s.Sum(nil))
}

// Returns hash: SHA256( SHA256( data ) )
// Where possible, using ShaHash() should be a bit faster
func Sha2Sum(b []byte) (out [32]byte) {
	ShaHash(b, out[:])
	return
}

func RimpHash(in []byte, out []byte) {
	sha := sha256.New()
	sha.Write(in)
	rim := ripemd160.New()
	rim.Write(sha.Sum(nil)[:])
	copy(out, rim.Sum(nil))
}
