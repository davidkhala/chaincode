package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/contract-api"
)

func main() {

	var chaincode = contract_api.NewChaincode(&SmartContract{}, &StupidContract{})

	contract_api.Start(chaincode)

}
