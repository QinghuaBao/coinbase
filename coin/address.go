package coin

import (
	"bytes"
	"math/big"
)

type StealthAddr struct {
	Version   byte
	Options   byte
	ScanKey   [33]byte
	SpendKeys [][33]byte
	Sigs      byte
	Prefix    []byte
}

type BtcAddr struct {
	Version  byte
	Hash160  [20]byte // For a stealth address: it's HASH160
	Checksum []byte   // Unused for a stealth address
	Pubkey   []byte   // Unused for a stealth address
	Enc58str string

	*StealthAddr // if this is not nil, means that this is a stealth address

	// This is used only by the client
	Extra struct {
		Label  string
		Wallet string
		Virgin bool
	}
}

func NewAddrFromPubkey(in []byte, ver byte) (a *BtcAddr) {
	a = new(BtcAddr)
	a.Pubkey = make([]byte, len(in))
	copy(a.Pubkey[:], in[:])
	a.Version = ver
	RimpHash(in, a.Hash160[:])
	return
}

func (a *StealthAddr) BytesNoPrefix() []byte {
	b := new(bytes.Buffer)
	b.WriteByte(a.Version)
	b.WriteByte(a.Options)
	b.Write(a.ScanKey[:])
	b.WriteByte(byte(len(a.SpendKeys)))
	for i := range a.SpendKeys {
		b.Write(a.SpendKeys[i][:])
	}
	b.WriteByte(a.Sigs)
	return b.Bytes()
}

func (a *StealthAddr) Bytes(checksum bool) []byte {
	b := new(bytes.Buffer)
	b.Write(a.BytesNoPrefix())
	b.Write(a.Prefix)
	if checksum {
		sh := Sha2Sum(b.Bytes())
		b.Write(sh[:4])
	}
	return b.Bytes()
}

func (a *StealthAddr) String() string {
	return Encodeb58(a.Bytes(true))
}

func (a *BtcAddr) String() string {
	if a.Enc58str == "" {
		if a.StealthAddr != nil {
			a.Enc58str = a.StealthAddr.String()
		} else {
			var ad [25]byte
			ad[0] = a.Version
			copy(ad[1:21], a.Hash160[:])
			if a.Checksum == nil {
				sh := Sha2Sum(ad[0:21])
				a.Checksum = make([]byte, 4)
				copy(a.Checksum, sh[:4])
			}
			copy(ad[21:25], a.Checksum[:])
			a.Enc58str = Encodeb58(ad[:])
		}
	}
	return a.Enc58str
}

var b58set []byte = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
var bn0 *big.Int = big.NewInt(0)
var bn58 *big.Int = big.NewInt(58)

func Encodeb58(a []byte) (s string) {
	idx := len(a)*138/100 + 1
	buf := make([]byte, idx)
	bn := new(big.Int).SetBytes(a)
	var mo *big.Int
	for bn.Cmp(bn0) != 0 {
		bn, mo = bn.DivMod(bn, bn58, new(big.Int))
		idx--
		buf[idx] = b58set[mo.Int64()]
	}
	for i := range a {
		if a[i] != 0 {
			break
		}
		idx--
		buf[idx] = b58set[0]
	}

	s = string(buf[idx:])

	return
}
