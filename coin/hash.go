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

// Returns hash: RIMP160( SHA256( data ) )
// Where possible, using RimpHash() should be a bit faster
func Rimp160AfterSha256(b []byte) (out [20]byte) {
	RimpHash(b, out[:])
	return
}

// This function is used to sign and verify messages using the bitcoin standard.
// The second paramater must point to a 32-bytes buffer, where hash will be stored.
func HashFromMessage(msg []byte, out []byte) {
	b := new(bytes.Buffer)
	WriteVlen(b, uint64(len(MessageMagic)))
	b.Write([]byte(MessageMagic))
	WriteVlen(b, uint64(len(msg)))
	b.Write(msg)
	ShaHash(b.Bytes(), out)
}

func WriteVlen(b io.Writer, var_len uint64) {
	if var_len < 0xfd {
		b.Write([]byte{byte(var_len)})
		return
	}
	if var_len < 0x10000 {
		b.Write([]byte{0xfd})
		binary.Write(b, binary.LittleEndian, uint16(var_len))
		return
	}
	if var_len < 0x100000000 {
		b.Write([]byte{0xfe})
		binary.Write(b, binary.LittleEndian, uint32(var_len))
		return
	}
	b.Write([]byte{0xff})
	binary.Write(b, binary.LittleEndian, var_len)
}
