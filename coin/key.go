package coin

import (
	"encoding/hex"
	"errors"

	"github.com/hyperledger/fabric/coinbase/lib-go-coin"
)

type PublicKey struct {
	secp256k1.XY
}

type Signature struct {
	secp256k1.Signature
	HashType byte
}


// Recoved public key form a signature
func (sig *Signature) RecoverPublicKey(msg []byte, recid int) (key *PublicKey) {
	key = new(PublicKey)
	if !secp256k1.RecoverPublicKey(sig.R.Bytes(), sig.S.Bytes(), msg, recid, &key.XY) {
		key = nil
	}
	return
}

func (sig *Signature) IsLowS() bool {
	return sig.S.Cmp(&secp256k1.TheCurve.HalfOrder.Int) < 1
}

// Returns serialized canoncal signature followed by a hash type
func (sig *Signature) Bytes() []byte {
	return append(sig.Signature.Bytes(), sig.HashType)
}
