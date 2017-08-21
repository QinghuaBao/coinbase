package coin

import "errors"

var (
	// ErrInvalidArgs is returned if there are some unused args or not enough args in params
	ErrInvalidArgs = errors.New("invalid args")

	// ErrInvalidFunction is returnd if chaincode interface get unsupported function name
	ErrInvalidFunction = errors.New("invalid function")

	// ErrInvalidTxKey returned if given key is invalid
	ErrInvalidTxKey = errors.New("invalid tx key")

	// ErrInvalidTX
	ErrInvalidTX = errors.New("transaction invalid")

	// ErrUnsupportedOperation returned if invoke or query using unsupported function name
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrMustCoinbase
	ErrMustCoinbase = errors.New("tx must be coinbase")

	// ErrCantCoinbase
	ErrCantCoinbase = errors.New("tx must not be coinbase")

	// ErrTxInOutNotBalance returned when txouts + fee != txins
	ErrTxInOutNotBalance = errors.New("tx in & out not balance")

	// ErrTxOutMoreThanTxIn
	ErrTxOutMoreThanTxIn = errors.New("tx out more than tx in")

	// ErrKeyNoData
	ErrKeyNoData = errors.New("state key found, but no data")

	// ErrCollisionTxOut
	ErrCollisionTxOut = errors.New("account has collision tx out")

	// ErrTxNoFounder
	ErrTxNoFounder = errors.New("tx has no founder")

	// ErrAccountNoTxOut
	ErrAccountNoTxOut = errors.New("account has no such tx out")

	// ErrAccountNotEnoughBalance
	ErrAccountNotEnoughBalance = errors.New("account has not enough balance")

	// ErrTxOutLock
	ErrTxOutLock = errors.New("tx out can be spend only after until time")

	// ErrAlreadyRegisterd
	ErrAlreadyRegisterd = errors.New("the addr has been registerd into coin")
	// ErrOutValueNegative
	ErrOutValueNegative = errors.New("txout value negative")
)

