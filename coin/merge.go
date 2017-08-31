package coin

import (
	"sort"

	"github.com/golang/protobuf/proto"
	//"github.com/hyperledger/fabric/coinbase/client"
	//"github.com/hyperledger/fabric/invoke_coinbase/proto"
)

func (coin *Hydruscoin) merge(store Store, account *Account) ([]byte, error) {
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

func NewTxIn(owner, prevHash string, prevIdx uint32) *TX_TXIN {
	return &TX_TXIN{
		SourceHash: prevHash,
		Ix:         prevIdx,
		Addr:       owner,
		Script:     "",
		Undefined:  "",
	}
}
