package main

import (
	"github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {

	contracts := []contractapi.ContractInterface{&SmartContract{}, &StupidContract{}}
	chaincode, err := contractapi.NewChaincode(contracts...)
	goutils.PanicError(err)

	err = chaincode.Start()

	goutils.PanicError(err)
}
