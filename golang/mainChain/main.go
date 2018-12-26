package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "mainChain"
)

type MainChaincode struct {
	CommonChaincode
}

func (t MainChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Init", fcn, params)
	if fcn != "" {
		t.PutState(fcn, []byte(params[0]))
	}
	return shim.Success(nil)
}

func (t MainChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
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
	t.Logger.Debug("response", responseBytes)
	return shim.Success(responseBytes)
}

func main() {
	var cc = MainChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
