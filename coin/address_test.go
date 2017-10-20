package coin

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	//"github.com/hyperledger/fabric/coinbase/secp256k1"
)

func TestNewAddrFromPubkey(t *testing.T) {
	//pubkey, privatekey := secp256k1.GenerateKeyPair()
	//pubkeystring := hex.EncodeToString(pubkey)
	//privatekeystring := hex.EncodeToString(privatekey)
	//fmt.Println(pubkeystring)
	//fmt.Println(privatekeystring)

	privatekey, err := hex.DecodeString("2C0D42397C4575E3DC0CD54599D3ECF342EE15DC33C5E876D8DC6AA7F3D280B0")
	pubstring := "BNf94GVnQ1XNecdDCVNBJhzhtUDGMYwWZrFWjAbTQwJLt8Jk0ye82OonfiaOaBYxQvLqE46sUPV04EOAmRluH1M="

	pubbyte, err := base64.StdEncoding.DecodeString(pubstring)
	fmt.Println("hex", byteToHexString(pubbyte))
	pubkey, err := hex.DecodeString(byteToHexString(pubbyte))
	//pubkey, err := hex.DecodeString("049CA4144418312C38B5189B37F4117322C4F75AE24EF2A1D7DADAA010AE93AC04B8F86CA18D439D5120EDACF33BD92B3CE14A4123853DCAD86A88375F8C015935")
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
