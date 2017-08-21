package coin

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
)

func stringToBigint(str string) (*big.Int, *big.Int) {
	strSplit := strings.Split(str, ":")

	r := new(big.Int)
	s := new(big.Int)

	r.SetString(strSplit[0], 10)
	s.SetString(strSplit[1], 10)

	return r, s
}

func signedmessage(tx *TX) ([]byte, error) {
	txsigned := new(TX)
	txsigned = tx

	script := tx.Txin[0].Script

	//	fmt.Println("txsigned txin", txsigned.Txin)
	//	fmt.Println("txsigned txout", txsigned.Txout)

	for index, _ := range txsigned.Txin {

		txsigned.Txin[index].Script = ""
	}

	//	fmt.Println("signedmessage ", txsigned)
	//fmt.Println(tx)

	txhash := TxHash(txsigned)
	signmessage, err := hex.DecodeString(txhash)
	if err != nil {
		logger.Errorf("hex.DecodeString error : %v", err)
		return signmessage, err
	}

	for index, _ := range txsigned.Txin {

		txsigned.Txin[index].Script = script
	}

	return signmessage, nil
}
func Verify(tx *TX) bool {
	verifymessage, err := signedmessage(tx)
	if err != nil {
		logger.Errorf("signedmessage error : %v", err)
		return false
	}

	if (len(tx.Txin) < 1) || (len(tx.Txout) < 1) {
		logger.Errorf("There is no txin or txout")
		return false
	}

	pubScr := tx.Txout[0].GetScriptPubKey()
	sigScript := tx.Txin[0].GetScript()
	logger.Debugf("pubScr: %v, verifymessage: %v, sigScript: %v", pubScr, base64.StdEncoding.EncodeToString(verifymessage), sigScript)
	cmd := exec.Command("java", "-jar", "Verify.jar", pubScr, base64.StdEncoding.EncodeToString(verifymessage), sigScript)
	out, err := cmd.Output()
	if err != nil {
		logger.Errorf("exec Command Error: %v", err)
		return false
	}
	fmt.Println(out)
	if strings.EqualFold(string(out), "false") {
		logger.Errorf("Verify Error: false")
		return false
	}

	return true

}
