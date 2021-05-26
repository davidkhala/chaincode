package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "stress"
)

var logger = shim.NewLogger(name)

type StressChaincode struct {
}

func (t StressChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Init")
	return shim.Success(nil)
}

func (t StressChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Invoke")
	return shim.Success(nil)
}

func main() {
	err := shim.Start(StressChaincode{})
	if err != nil {
		panic(err)
	}
}
