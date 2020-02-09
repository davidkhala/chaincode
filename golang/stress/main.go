package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	logger "github.com/sirupsen/logrus"
)

const (
	name = "stress"
)

type StressChaincode struct {
}

func (t *StressChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Init")
	return shim.Success(nil)
}

func (t *StressChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("Invoke")
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(StressChaincode))
	if err != nil {
		panic(err)
	}
}
