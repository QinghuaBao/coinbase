package coin

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const coinInfoKey = "HydruscoinInfo"

// Key represents the key for a transaction in storage. It has both a
// hash and index
type Key struct {
	TxHashAsHex string
	TxIndex     uint32
}

func (k *Key) String() string {
	return fmt.Sprintf("%s:%d", k.TxHashAsHex, k.TxIndex)
}

// parseKey parse key string into Key object, return error if something invalid happened
func parseKey(keyStr string) (*Key, error) {
	if !strings.Contains(keyStr, ":") {
		return nil, ErrInvalidTxKey
	}

	subKeys := strings.Split(keyStr, ":")
	if len(subKeys) != 2 {
		return nil, ErrInvalidTxKey
	}

	txIdx, err := strconv.ParseUint(subKeys[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &Key{TxHashAsHex: subKeys[0], TxIndex: uint32(txIdx)}, nil
}

func generateAccountKey(addr string) string {
	return fmt.Sprintf("account_addr_%s", addr)
}

// Store interface describes the storage used by this chaincode. The interface
// was created so either the state database store can be used or a in memory
// store can be used for unit testing.
type Store interface {
	GetTx(string) (*TX, bool, error)
	PutTx(*TX) error
	InitCoinInfo() error
	GetCoinInfo() (*HydruscoinInfo, error)
	PutCoinInfo(*HydruscoinInfo) error
	GetAccount(string) (*Account, error)
	PutAccount(*Account) error
	GetAccountTxout(addr string, num uint32) (*Account, error)
}

// Store struct uses a chaincode stub for state access
type ChaincodeStore struct {
	stub shim.ChaincodeStubInterface
}

// MakeChaincodeStore returns a store for storing keys in the state
func MakeChaincodeStore(stub shim.ChaincodeStubInterface) Store {
	store := &ChaincodeStore{}
	store.stub = stub
	return store
}

// GetTx returns a transaction for the given hash
func (s *ChaincodeStore) GetTx(key string) (*TX, bool, error) {
	data, err := s.stub.GetState(key)
	if err != nil {
		return nil, false, fmt.Errorf("Error getting state from stub:  %s", err)
	}
	if data == nil || len(data) == 0 {
		return nil, false, nil
	}

	tx, err := ParseTXBytes(data)
	if err != nil {
		return nil, false, err
	}

	return tx, true, nil
}

// PutTx adds a transaction to the state with the hash as a key
func (s *ChaincodeStore) PutTx(tx *TX) error {
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		return err
	}

	return s.stub.PutState(TxHash(tx), txBytes)
}

func (s *ChaincodeStore) InitCoinInfo() error {
	coinInfo := &HydruscoinInfo{
		CoinTotal:    0,
		AccountTotal: 0,
		TxoutTotal:   0,
		TxTotal:      0,
		Session: &HydruscoinInfo_POVSession{
			TxCount:               0,
			CurrentAlpha:          0.7,
			CurrentTotalIncentive: 0,
			Threshold:             100,
		},
		Placeholder: "placeholder",
	}

	return s.PutCoinInfo(coinInfo)
}

func (s *ChaincodeStore) GetCoinInfo() (*HydruscoinInfo, error) {
	data, err := s.stub.GetState(coinInfoKey)
	if err != nil {
		return nil, err
	}

	if data == nil || len(data) == 0 {
		return nil, ErrKeyNoData
	}

	coinfo, err := ParseHydruscoinInfoBytes(data)
	if err != nil {
		return nil, err
	}

	return coinfo, nil
}

func (s *ChaincodeStore) PutCoinInfo(coinfo *HydruscoinInfo) error {
	coinBytes, err := proto.Marshal(coinfo)
	if err != nil {
		return err
	}

	if err := s.stub.PutState(coinInfoKey, coinBytes); err != nil {
		return err
	}

	return nil
}

// GetAccount returns account from world states
func (s *ChaincodeStore) GetAccount(addr string) (*Account, error) {
	if addr == "" {
		return nil, errors.New("empty addr")
	}
	key := generateAccountKey(addr)
	data, err := s.stub.GetState(key)
	if err != nil {
		return nil, err
	}

	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("no account found")
	}

	accountslice := new(AccountSlice)
	if err := proto.Unmarshal(data, accountslice); err != nil {
		return nil, err
	}

	account := ParseAccount(accountslice)

	return account, nil
}

// PutAccount update or insert account into world states
func (s *ChaincodeStore) PutAccount(account *Account) error {
	key := generateAccountKey(account.Addr)

	accountslice := GenerateAccount(account)

	aBytes, err := proto.Marshal(accountslice)
	if err != nil {
		return err
	}

	return s.stub.PutState(key, aBytes)
}

//added
// GetAccountTxout returns account from world states
func (s *ChaincodeStore) GetAccountTxout(addr string, num uint32) (*Account, error) {
	if addr == "" || num == 0 {
		return nil, errors.New("empty addr or zero num")
	}
	key := generateAccountKey(addr)
	data, err := s.stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data) == 0 {
		return nil, fmt.Errorf("no account found")
	}
	accountslice := new(AccountSlice)
	tmpaccslice := new(AccountSlice)
	if err := proto.Unmarshal(data, tmpaccslice); err != nil {
		return nil, err
	}
	if num > uint32(len(tmpaccslice.Txoutmap)) {
		logger.Debugf("query num[%d] exceed count: %+v")
		num = uint32(len(tmpaccslice.Txoutmap))
	}
	accountslice.Balance = tmpaccslice.Balance
	accountslice.Addr = tmpaccslice.Addr
	var i uint32
	accountslice.Txoutmap = make([]*TxoutMap, num)
	for i = 0; i < num; i++ {
		accountslice.Txoutmap[i] = tmpaccslice.Txoutmap[i]
	}
	account := ParseAccount(accountslice)
	return account, nil
}
