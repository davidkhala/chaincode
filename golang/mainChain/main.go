package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "mainChain"
)

type MainChaincode struct {
	CommonChaincode
}

func (t MainChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Init")

	return shim.Success(nil)
}

func (t MainChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Invoke")
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Debug("fcn", fcn, "params", params)
	var key = params[0]
	var responseBytes []byte
	switch fcn {
	case "put":
		t.PutState(key, []byte(params[1]))
	case "get":
		responseBytes = t.GetState(key)
	}

	return shim.Success(responseBytes)
}

func main() {
	var cc = MainChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
