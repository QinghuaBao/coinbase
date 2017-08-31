package coin

import (
	"encoding/base64"
	//	"encoding/hex"
	//	"strings"
	//"time"

	"github.com/golang/protobuf/proto"
	//"github.com/hyperledger/fabric/coinbase/sql"
)

func (coin *Hydruscoin) coinbase(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		logger.Errorf("Decoding base64 error: %v\n", err)
		return nil, err
	}

	tx, err := ParseTXBytes(txData)
	if err != nil {
		logger.Errorf("Unmarshal tx bytes error: %v\n", err)
		return nil, err
	}
	logger.Debugf("tx: %v", tx)

	//if len(tx.Txin) == 0 {
	//	return nil, ErrInvalidTX
	//}
	txhash := TxHash(tx)
	execResult := &ExecResult{}

	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	// Loop through outputs
	for index, output := range tx.Txout {
		//verfiy address
		flag := verifyAddr(output.GetScriptPubKey(), output.Addr, tx.Version)
		if !flag {
			return nil, ErrInvalidTX
		}

		//		pubkey, err := hex.DecodeString(output.ScriptPubKey)
		//		if err != nil {
		//			logger.Errorf("Error hex.DecodeString : %v", err)
		//			return nil, err
		//		}
		//		addr := NewAddrFromPubkey(pubkey, byte(tx.Version))
		//		//		logger.Debugf("addr: %v, version: %s", pubkey, tx.Version)
		//		logger.Debugf("newaddrfrompubkey: %v, addr: %v", addr, output.Addr)
		//		if !strings.EqualFold(addr.String(), output.Addr) {
		//			logger.Errorf("genaddr: %v, addr: %v", addr.String(), output.Addr)
		//			return nil, ErrInvalidTX
		//		}

		if output.Addr == "" {
			logger.Errorf("output.Addr is null")
			return nil, ErrInvalidTX
		}

		outerAccount, err := store.GetAccount(output.Addr)
		if err != nil {
			logger.Warningf("account[%s] is not existed, creating one...", output.Addr)

			outerAccount = new(Account)
			outerAccount.Addr = output.Addr
			outerAccount.Txouts = make(map[string]*TX_TXOUT)

			coinInfo.AccountTotal += 1
		}
		if outerAccount.Txouts == nil || len(outerAccount.Txouts) == 0 {
			outerAccount.Txouts = make(map[string]*TX_TXOUT)
		}
		//		logger.Debugf("output.Addr == ")
		currKey := &Key{TxHashAsHex: txhash, TxIndex: uint32(index)}
		if _, ok := outerAccount.Txouts[currKey.String()]; ok {
			return nil, ErrCollisionTxOut
		}
		if output.Value < 0 {
			return nil, ErrOutValueNegative
		}
		// store tx out into account
		outerAccount.Txouts[currKey.String()] = output
		outerAccount.Balance += output.Value
		if err := store.PutAccount(outerAccount); err != nil {
			logger.Errorf("Error update account: %v, account info: %+v", err, outerAccount)
			return nil, err
		}
		logger.Debugf("put tx output %s:%v", currKey.String(), output)

		// change coin info
		coinInfo.CoinTotal += output.Value
		coinInfo.TxoutTotal += 1
		execResult.SumCurrentOutputs += output.Value
	}

	//sava tx information to mysql and blockchain
	if err := store.PutTx(tx); err != nil {
		logger.Errorf("put tx error: %v", err)
		return nil, err
	}
	logger.Debug("put tx into world state")

	//merge
	logger.Debugf("11")
	account, err := store.GetAccount(tx.Txout[0].Addr)
	if err != nil {
		return nil, err
	}

	logger.Debugf("22s")
	if len(account.Txouts) >= 1000 {
		if _, err := coin.merge(store, account); err != nil {
			return nil, err
		}
	}

	//err = sql.InsertTran(txhash, time.Now().UTC().Unix(), "coinbase", tx.Txout[0].Addr, execResult.SumCurrentOutputs, tx.Version, tx.Founder, execResult.SumCurrentOutputs, "coinbase")
	//if err != nil {
	//	logger.Errorf("insert sql error : %V", err)
	//	return nil, err
	//}
	logger.Debugf("put tx into mysql")

	// tx total counter
	coinInfo.TxTotal += 1
	if err := store.PutCoinInfo(coinInfo); err != nil {
		logger.Errorf("Error put coin info: %v", err)
		return nil, err
	}

	logger.Debugf("coinbase execute result: %+v", execResult)
	return proto.Marshal(execResult)
}
