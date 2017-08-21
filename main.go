package main

import (
	"github.com/hyperledger/fabric/coinbase/coin"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/op/go-logging"
)

func main() {
	//shim.SetLoggingLevel(shim.LoggingLevel(logging.ERROR))
	if err := shim.Start(&coin.Hydruscoin{}); err != nil {
		panic(err)
	}
}
