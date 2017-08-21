package coin

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/coinbase/secp256k1"
)

func TestNewAddrFromPubkey(t *testing.T) {
	pubkey, privatekey := secp256k1.GenerateKeyPair()
	//pubkeystring := hex.EncodeToString(pubkey)
	//privatekeystring := hex.EncodeToString(privatekey)
	//fmt.Println(pubkeystring)
	//fmt.Println(privatekeystring)

	privatekey, err := hex.DecodeString("2C0D42397C4575E3DC0CD54599D3ECF342EE15DC33C5E876D8DC6AA7F3D280B0")
	pubstring := "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAESO/MSgx3EaaazRFb4SDrfKBHYFx9yF9/gTIqnnOz+kDZUWf8M7vYHMHg4Tfh9++vKysaOzi5BerMFHsS8kdWZA=="

	pubbyte, err := base64.StdEncoding.DecodeString(pubstring)
	fmt.Println("hex", byteToHexString(pubbyte))
	pubkey, err = hex.DecodeString(byteToHexString(pubbyte))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pubkey, privatekey)

	var version int
	version = 1

	addr := NewAddrFromPubkey(pubkey, byte(version))
	//fmt.Println(addr)
	fmt.Println(addr.String())
}
