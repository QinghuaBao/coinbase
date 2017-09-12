package coin

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/flogging"
	"github.com/op/go-logging"
)

var (
	logger = logging.MustGetLogger("foam")
)

// Hydruscoin
type Hydruscoin struct{}

// Init deploy chaincode into vp
func (coin *Hydruscoin) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	flogging.LoggingInit("hydruscoin")
	if function != "deploy" {
		return nil, ErrInvalidFunction
	}

	// construct a new store
	store := MakeChaincodeStore(stub)

	// deploy hydruscoin chaincode only need to set coin stater
	if err := store.InitCoinInfo(); err != nil {
		return nil, err
	}

	logger.Debug("deploy Hydruscoin successfully")
	return nil, nil
}

// Invoke function
const (
	IF_REGISTER string = "invoke_register"
	IF_COINBASE string = "invoke_coinbase"
	IF_TRANSFER string = "invoke_transfer"
	IF_POV      string = "invoke_transfer_pov"
)

// Invoke
func (coin *Hydruscoin) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case IF_REGISTER:
		return coin.registerAccount(store, args)
	case IF_COINBASE:
		return coin.coinbase(store, args)
	case IF_TRANSFER:
		return coin.transfer(store, args)
	case IF_POV:
		return coin.pov_transfer(store, args)
	default:
		return nil, ErrUnsupportedOperation
	}
}

// Query function
const (
	QF_ADDRS       = "query_addrs"
	QF_TX          = "query_tx"
	QF_COIN        = "query_coin"
	QF_BALANCE     = "query_balance"
	QF_TXOUT_COUNT = "query_txout_count"
	QF_ADDRS_TXOUT = "query_addrs_txout"
)

// Query
func (coin *Hydruscoin) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	// construct a new store
	store := MakeChaincodeStore(stub)

	switch function {
	case QF_ADDRS:
		return coin.queryAddrs(store, args)
	case QF_TX:
		return coin.queryTx(store, args)
	case QF_COIN:
		return coin.queryCoin(store, args)
	case QF_BALANCE:
		return coin.queryBalance(store, args)
	case QF_TXOUT_COUNT:
		return coin.queryTxoutCount(store, args)
	case QF_ADDRS_TXOUT:
		return coin.queryAddrsTxout(store, args)
	default:
		return nil, ErrUnsupportedOperation
	}
}
