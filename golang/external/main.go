package main

import (
	"github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func main() {

	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	goutils.PanicError(err)
	// not chaincode as a service

	err = shim.Start(chaincode) // mandate CORE_CHAINCODE_ID_NAME

	goutils.PanicError(err)
}
