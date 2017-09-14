package coin

import (
	"encoding/base64"

	//	"encoding/binary"
	//	"math"
	"strconv"

	"github.com/golang/protobuf/proto"
)

func (coin *Hydruscoin) queryAddrs(store Store, args []string) ([]byte, error) {
	results := &QueryAddrResults{
		Accounts: make(map[string]*Account),
	}

	for _, addr := range args {
		account, err := store.GetAccount(addr)
		if err != nil {
			logger.Errorf("store.GetAccount() return error: %v", err)
			continue
		}

		results.Accounts[addr] = account
		logger.Debugf("query addr[%s] account: %+v", addr, account)
	}

	protobyte, err := proto.Marshal(results)
	if err != nil {
		logger.Errorf("result marshal error: %v", err)
		return nil, err
	}

	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(protobyte)))
	//var protobytebase64 [0:len(protobyte)]byte
	base64.StdEncoding.Encode(protobytebase64, protobyte)
	//logger.Debugf("protobytebase64: %v", protobytebase64)

	return protobytebase64, err
}

//改数组
func (coin *Hydruscoin) queryTx(store Store, args []string) ([]byte, error) {
	if len(args) != 1 || args[0] == "" {
		return nil, ErrInvalidArgs
	}

	tx, _, err := store.GetTx(args[0])
	if err != nil {
		logger.Errorf("get tx info error: %v", err)
		return nil, err
	}
	logger.Debugf("query tx: %+v", tx)

	protobyte, err := proto.Marshal(tx)
	if err != nil {
		logger.Debugf("tx marshal error: %v", err)
		return nil, err
	}

	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(protobyte)))
	base64.StdEncoding.Encode(protobytebase64, protobyte)

	return protobytebase64, err
}

func (coin *Hydruscoin) queryCoin(store Store, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, ErrInvalidArgs
	}

	coinInfo, err := store.GetCoinInfo()
	if err != nil {
		logger.Errorf("Error get coin info: %v", err)
		return nil, err
	}

	logger.Debugf("query lepuscoin info: %+v", coinInfo)

	protobyte, err := proto.Marshal(coinInfo)
	if err != nil {
		logger.Debugf("tx marshal error: %v", err)
		return nil, err
	}

	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(protobyte)))
	base64.StdEncoding.Encode(protobytebase64, protobyte)
	return protobytebase64, err
}

//added

func (coin *Hydruscoin) queryBalance(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}
	addr := args[0]
	account, err := store.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	logger.Debugf("query balance[%s] account: %+v", addr, account)
	//binary.BigEndian.PutUint32(bytes, uint32(account.Balance))
	s := strconv.FormatInt(account.Balance, 10)
	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(protobytebase64, []byte(s))
	return protobytebase64, err
}

func (coin *Hydruscoin) queryTxoutCount(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}
	addr := args[0]
	account, err := store.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	logger.Debugf("query txout count[%s] account: %+v", addr, account)
	//	var bytes = make([]byte, 4)
	s := strconv.Itoa(len(account.Txouts))
	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(protobytebase64, []byte(s))
	return protobytebase64, err
}

func (coin *Hydruscoin) queryAddrsTxout(store Store, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, ErrInvalidArgs
	}
	results := &QueryAddrResults{
		Accounts: make(map[string]*Account),
	}
	addr := args[0]
	num, err := strconv.Atoi(args[1])
	if err != nil {
		logger.Errorf("unable to convert num: %v", err)
		return nil, err
	}
	account, err := store.GetAccountTxout(string(addr), uint32(num))
	if err != nil {
		logger.Errorf("store.GetAccount() return error: %v", err)
		return nil, err
	}
	results.Accounts[addr] = account
	logger.Debugf("query addr[%s] account: %+v", addr, account)
	protobyte, err := proto.Marshal(results)
	if err != nil {
		logger.Errorf("result marshal error: %v", err)
		return nil, err
	}
	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(protobyte)))
	base64.StdEncoding.Encode(protobytebase64, protobyte)
	return protobytebase64, err
}

func (coin *Hydruscoin) queryTest(store Store, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, ErrInvalidArgs
	}

	coinTest, err := store.GetTest()
	if err != nil {
		logger.Errorf("Error get coin test: %v", err)
		return nil, err
	}

	logger.Debugf("query lepuscoin info: %+v", coinTest)

	protobyte, err := proto.Marshal(coinTest)
	if err != nil {
		logger.Debugf("tx marshal error: %v", err)
		return nil, err
	}

	protobytebase64 := make([]byte, base64.StdEncoding.EncodedLen(len(protobyte)))
	base64.StdEncoding.Encode(protobytebase64, protobyte)
	return protobytebase64, err
}
