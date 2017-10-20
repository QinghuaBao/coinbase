package btcd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcd/btcec"
	"github.com/btcd/chaincfg/chainhash"
)

func TestEncryption(t *testing.T) {
	// Decode a hex-encoded private key.
	s := "22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c"
	pkBytes, err := hex.DecodeString(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	// Sign a message using the private key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, err := privKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		return
	}
	//x1 := signature.Serialize();

	//ansss := base64.StdEncoding.EncodeToString(x1)

	//fmt.Printf("---base64: %s\n\n", ansss)
	fmt.Printf("r:%s, s:%s\n", signature.R.String(), signature.S.String())

	// Verify the signature for the message using the public key.
	verified := signature.Verify(messageHash, pubKey)
	fmt.Printf("Signature Verified? %v\n", verified)

	//start
	fmt.Print("verify..\n")
	sig := "MEUCIQDYSSLig/iYjPWMZHPBppvJlXaO2dUmGVY6gC5q0v8gIwIgfTkjNm3uRIJ+x5DMMr2H+3lBuxvuzaXth9tjRpbPJHk="

	sigB, _ := base64.StdEncoding.DecodeString(sig)
	signature0, _ := btcec.ParseSignature(sigB, btcec.S256())

	fmt.Printf("r:%s, s:%s\n", signature0.R.String(), signature0.S.String())

	//message0 := "test message"

	pubKeyStr := "AmLWWi4u6d2uOw8VPytfNMDdtKQYB5JEWFVIYHc8QeMy"
	pubKeyB, _ := base64.StdEncoding.DecodeString(pubKeyStr)

	pubkey0, _ := btcec.ParsePubKey(pubKeyB, btcec.S256())
	verified0 := signature0.Verify([]byte(messageHash), pubkey0)

	fmt.Printf("--Signature Verified? %v\n", verified0)
	// Output:
	// Serialized Signature: 304402201008e236fa8cd0f25df4482dddbb622e8a8b26ef0ba731719458de3ccd93805b022032f8ebe514ba5f672466eba334639282616bb3c2f0ab09998037513d1f9e3d6d
	// Signature Verified? true
}

//package main
//
//import (
//	"encoding/hex"
//	"fmt"
//
//	"mycode/coinbase/btcd/btcec"
//	"mycode/coinbase/btcd/chaincfg/chainhash"
//	"encoding/base64"
//)

//func main() {
//	// Decode a hex-encoded private key.
//	s := "22a47fa09a223f2aa079edf85a7c2d4f87" +
//		"20ee63e502ee2869afab7de234b80c";
//	pkBytes, err := hex.DecodeString(s)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
//	// Sign a message using the private key.
//	message := "test message"
//	messageHash := chainhash.DoubleHashB([]byte(message))
//	signature, err := privKey.Sign(messageHash)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//x1 := signature.Serialize();
//
//	//ansss := base64.StdEncoding.EncodeToString(x1)
//
//	//fmt.Printf("---base64: %s\n\n", ansss)
//	fmt.Printf("r:%s, s:%s\n",signature.R.String(), signature.S.String())
//
//	// Verify the signature for the message using the public key.
//	verified := signature.Verify(messageHash, pubKey)
//	fmt.Printf("Signature Verified? %v\n", verified)
//
//
//	//start
//	fmt.Print("verify..\n")
//	sig := "MEUCIQDYSSLig/iYjPWMZHPBppvJlXaO2dUmGVY6gC5q0v8gIwIgfTkjNm3uRIJ+x5DMMr2H+3lBuxvuzaXth9tjRpbPJHk="
//
//	sigB, _ := base64.StdEncoding.DecodeString(sig)
//	signature0, _ := btcec.ParseSignature(sigB, btcec.S256())
//
//	fmt.Printf("r:%s, s:%s\n",signature0.R.String(), signature0.S.String())
//
//	//message0 := "test message"
//
//	pubKeyStr := "AmLWWi4u6d2uOw8VPytfNMDdtKQYB5JEWFVIYHc8QeMy"
//	pubKeyB, _ := base64.StdEncoding.DecodeString(pubKeyStr)
//
//	pubkey0,_ := btcec.ParsePubKey(pubKeyB, btcec.S256())
//	verified0 := signature0.Verify([]byte(messageHash), pubkey0)
//
//	fmt.Printf("--Signature Verified? %v\n", verified0)
//	// Output:
//	// Serialized Signature: 304402201008e236fa8cd0f25df4482dddbb622e8a8b26ef0ba731719458de3ccd93805b022032f8ebe514ba5f672466eba334639282616bb3c2f0ab09998037513d1f9e3d6d
//	// Signature Verified? true
//}
