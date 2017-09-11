package coin

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
	//"github.com/hyperledger/fabric/coinbase/sql"

	"github.com/golang/protobuf/proto"
	"math"
)

func (coin *Hydruscoin) transfer(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	// parse tx
	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		logger.Errorf("Error decode tx bytes: %v", err)
		return nil, err
	}

	tx, err := ParseTXBytes(txData)
	if err != nil {
		return nil, err
	}
	logger.Debugf("tx: %v", tx)

	//只支持一种来源的txin
	if len(tx.Txin) == 0 {
		return nil, ErrInvalidTX
	}
	//verify script
	logger.Debugf("hhhh")
	ok := Verify(tx)
	if ok != true {
		logger.Errorf("Verify Error ")
		return nil, errors.New("Verify Error")
	}

	// coin stat
	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	execResult := &ExecResult{}
	txHash := TxHash(tx)
	if tx.Founder == "" {
		return nil, ErrTxNoFounder
	}
	logger.Debugf("tx: %v", tx)
	//founderAccount, err := store.GetAccount(tx.Founder)
	//if err != nil {
	//	return nil, ErrTxNoFounder
	//}

	for _, ti := range tx.Txin {
		prevTxHash := ti.SourceHash
		prevOutputIx := ti.Ix
		ownerAddr := ti.Addr
		keyToPrevOutput := &Key{TxHashAsHex: prevTxHash, TxIndex: prevOutputIx}

		ownerAccount, err := store.GetAccount(ownerAddr)
		if err != nil {
			return nil, err
		}
		txout, ok := ownerAccount.Txouts[keyToPrevOutput.String()]
		if !ok {
			return nil, ErrAccountNoTxOut
		}

		// can spend?
		if txout.Until > 0 {
			untilTime := time.Unix(txout.Until, 0).UTC()
			if untilTime.After(time.Now().UTC()) {
				return nil, ErrTxOutLock
			}
		}

		if ownerAccount.Balance < txout.Value {
			return nil, ErrAccountNotEnoughBalance
		}
		ownerAccount.Balance -= txout.Value
		delete(ownerAccount.Txouts, keyToPrevOutput.String())

		// save owner account
		if err := store.PutAccount(ownerAccount); err != nil {
			return nil, err
		}

		// coin stat
		coinInfo.TxoutTotal -= 1
		execResult.SumPriorOutputs += txout.Value
	}

	// save founder account
	//if err := store.PutAccount(founderAccount); err != nil {
	//	return nil, err
	//}

	incentives := make(map[string]*TX_TXOUT)
	for idx, to := range tx.Txout {
		flag := verifyAddr(to.ScriptPubKey, tx.Txin[0].GetAddr(), tx.Version)
		if !flag {
			return nil, ErrInvalidTX
		}

		account, err := store.GetAccount(to.Addr)
		if err != nil {
			logger.Warningf("get account[%s] doesnt exist, creating one...", to.Addr)

			account = new(Account)
			account.Txouts = make(map[string]*TX_TXOUT)
			account.Addr = to.Addr

			coinInfo.AccountTotal += 1
		}
		if account.Txouts == nil || len(account.Txouts) == 0 {
			account.Txouts = make(map[string]*TX_TXOUT)
		}

		outKey := &Key{TxHashAsHex: txHash, TxIndex: uint32(idx)}
		if _, ok := account.Txouts[outKey.String()]; ok {
			return nil, ErrCollisionTxOut
		}
		if to.Value < 0 {
			return nil, ErrOutValueNegative
		}
		account.Balance += to.Value
		account.Txouts[outKey.String()] = to

		incentiveTxout, ok := incentives[to.GetAddr()]
		if ok != true {
			incentiveTxout = &TX_TXOUT{
				Addr:         to.GetAddr(),
				ScriptPubKey: to.ScriptPubKey,
				Value:        0,
				Until:        -1,
			}
			incentives[to.GetAddr()] = incentiveTxout
		}
		deltaIncentive := int64(math.Ceil(float64(coinInfo.Session.CurrentAlpha*float32(to.Value) + 0.5)))
		incentiveTxout.Value += deltaIncentive
		coinInfo.Session.CurrentTotalIncentive += deltaIncentive

		//save account
		if err := store.PutAccount(account); err != nil {
			return nil, err
		}

		// coin stat
		coinInfo.TxoutTotal += 1
		execResult.SumCurrentOutputs += to.Value
	}

	// current outputs must less than prior outputs
	if execResult.SumCurrentOutputs > execResult.SumPriorOutputs {
		return nil, ErrTxOutMoreThanTxIn
	}

	//	when courrent outputs more than prior, change
	if execResult.SumPriorOutputs > execResult.SumCurrentOutputs {
		ownerAccount, err := store.GetAccount(tx.Txin[0].Addr)
		if err != nil {
			return nil, err
		}
		if ownerAccount.Txouts == nil || len(ownerAccount.Txouts) == 0 {
			ownerAccount.Txouts = make(map[string]*TX_TXOUT)
		}
		changekey := &Key{TxHashAsHex: txHash, TxIndex: 0}
		if _, ok := ownerAccount.Txouts[changekey.String()]; ok {
			return nil, ErrCollisionTxOut
		}
		changeTxouts := new(TX_TXOUT)
		changeTxouts.Addr = ownerAccount.Addr
		changeTxouts.ScriptPubKey = ""
		changeTxouts.Until = -1
		changeTxouts.Value = execResult.SumPriorOutputs - execResult.SumCurrentOutputs
		changeTxouts.Undefined = ""

		ownerAccount.Txouts[changekey.String()] = changeTxouts
		ownerAccount.Balance += changeTxouts.Value
		coinInfo.TxoutTotal += 1

		// save owner account
		if err := store.PutAccount(ownerAccount); err != nil {
			return nil, err
		}
	}
	//merge
	account, err := store.GetAccount(tx.Txout[0].Addr)
	if err != nil {
		return nil, err
	}

	if len(account.Txouts) >= 1000 {
		if _, err := coin.merge(store, account); err != nil {
			return nil, err
		}
	}

	//if execResult.SumCurrentOutputs != execResult.SumPriorOutputs {
	//	return nil, ErrTxInOutNotBalance
	//}
	//sava tx information to mysql and blockchain
	if err := store.PutTx(tx); err != nil {
		logger.Errorf("put tx error: %v", err)
		return nil, err
	}
	logger.Debug("put tx into world state")

	//err = sql.InsertTran(txHash, time.Now().UTC().Unix(), tx.Txin[0].Addr, tx.Txout[0].Addr, execResult.SumCurrentOutputs, tx.Version, tx.Founder, execResult.SumPriorOutputs-execResult.SumCurrentOutputs, "transfer")
	//if err != nil {
	//	logger.Errorf("insert sql error : %V", err)
	//	return nil, err
	//}

	response, err := doIncentive(store, &incentives, tx.Version, coin)
	if err != nil {
		return response, err
	}
	if coinInfo.Session.CurrentTotalIncentive >= INCENT_THREADSHOLD {
		updatePovSession(store, coinInfo.Session)
	}
	logger.Debugf("put tx into mysql")

	// save coin stat
	coinInfo.TxTotal += 1
	if err := store.PutCoinInfo(coinInfo); err != nil {
		return nil, err
	}

	return proto.Marshal(execResult)
}

func byteToHexString(byteArray []byte) string {
	result := ""
	for i := 0; i < len(byteArray); i++ {
		hex := strconv.FormatInt(int64(byteArray[i]&0xFF), 16)
		if len(hex) == 1 {
			hex = "0" + hex
		}
		result += hex
	}
	return strings.ToUpper(result)
}

func verifyAddr(str string, addr string, version uint64) bool {
	pubbyte, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return false
	}
	pubkey, err := hex.DecodeString(byteToHexString(pubbyte))
	if err != nil {
		return false
	}

	addrObtain := NewAddrFromPubkey(pubkey, byte(version))
	logger.Debugf("newaddrfrompubkey : %v, addr: %v", addrObtain, addr)

	if !strings.EqualFold(addrObtain.String(), addr) {
		logger.Errorf("%v", ErrInvalidTX)
		return false
	}
	return true
}
