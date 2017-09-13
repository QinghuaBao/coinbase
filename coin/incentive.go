package coin

import (
	"github.com/golang/protobuf/proto"
	"encoding/base64"
	"sort"
)



//Dynamic adjustment
func updatePovSession(store Store, session *HydruscoinInfo_POVSession, timestamp int64) {
	logger.Debug("Adjusting POV session parameters")
	session.CurrentAlpha = (timestamp - session.TxCount)*session.CurrentAlpha / (INCENT_T0*100)
	if session.CurrentAlpha > 70 || session.CurrentAlpha <= 0{
		session.CurrentAlpha = 70
	}
	session.CurrentTotalIncentive = 0
	session.TxCount = timestamp
}

//timestamp make txout cant repeat
func doIncentive(store Store, incentives *map[string]*TX_TXOUT, version uint64, coin *Hydruscoin, timestamp int64) ([]byte, error) {
	logger.Debug("Doing Incentives")

	tx := &TX{Txout: make([]*TX_TXOUT, len(*incentives)), Version: version, Timestamp: timestamp, Founder: "blockchain"}

	i := 0
	for _, val := range *incentives {
		tx.Txout[i] = val
		i++
	}

	arg, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}
	base64 := base64.StdEncoding.EncodeToString(arg)
	return coinbase(store, []string{base64})
}

func coinbase(store Store, args []string) ([]byte, error) {
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
		//flag := verifyAddr(output.GetScriptPubKey(), output.Addr, tx.Version)
		//if !flag {
		//	return nil, ErrInvalidTX
		//}

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
		if _, err := merge(store, account); err != nil {
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

func merge(store Store, account *Account) ([]byte, error) {
	var value int64
	value = 0

	tx := new(TX)
	tx.Founder = "blockchain-foam"
	tx.Timestamp = 1060762481
	tx.Version = 1
	logger.Debugf("1")
	//固定到数组
	tx.Txin = make([]*TX_TXIN, len(account.Txouts))
	temp := make([]string, len(account.Txouts))
	i := 0
	for key, txout := range account.Txouts {
		temp[i] = key
		value += txout.Value
		i++
	}
	sort.Strings(temp)
	for j, key := range temp {
		preKey, err := parseKey(key)
		if err != nil {
			return nil, err
		}
		tx.Txin[j] = NewTxIn(account.Addr, preKey.TxHashAsHex, preKey.TxIndex)
	}

	tx.Txout = make([]*TX_TXOUT, 1)
	//tx.Txout[0] = NewTxOut(value, account.Addr, -1, "")
	tx.Txout[0] = &TX_TXOUT{
		Value:        value,
		Addr:         account.Addr,
		Until:        -1,
		ScriptPubKey: "",
		Undefined:    "",
	}

	txhash := TxHash(tx)
	logger.Debugf("2")
	execResult := &ExecResult{}
	execResult.SumPriorOutputs = value

	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	coinInfo.TxoutTotal -= int64(len(account.Txouts))
	account.Balance -= value
	account.Txouts = make(map[string]*TX_TXOUT)
	if err := store.PutAccount(account); err != nil {
		logger.Errorf("Error update account: %v, account info: %+v", err, account)
		return nil, err
	}

	logger.Debugf("3")
	// Loop through outputs
	for index, output := range tx.Txout {

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

		// store tx out into account
		outerAccount.Txouts[currKey.String()] = output
		outerAccount.Balance += output.Value
		if err := store.PutAccount(outerAccount); err != nil {
			logger.Errorf("Error update account: %v, account info: %+v", err, outerAccount)
			return nil, err
		}
		logger.Debugf("put tx output %s:%v", currKey.String(), output)

		// change coin info
		//coinInfo.CoinTotal += output.Value
		coinInfo.TxoutTotal += 1
		execResult.SumCurrentOutputs += output.Value
	}
	logger.Debugf("4")
	//sava tx information to mysql and blockchain
	if err := store.PutTx(tx); err != nil {
		logger.Errorf("put tx error: %v", err)
		return nil, err
	}
	logger.Debug("put tx into world state")

	// tx total counter
	coinInfo.TxTotal += 1
	if err := store.PutCoinInfo(coinInfo); err != nil {
		logger.Errorf("Error put coin info: %v", err)
		return nil, err
	}
	logger.Debugf("5")

	logger.Debugf("merge execute result: %+v", execResult)
	return proto.Marshal(execResult)
}
