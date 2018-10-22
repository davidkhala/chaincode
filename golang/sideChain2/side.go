package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name   = "sideChaincode2"
	mainCC = "mainChain"
)

type SideChaincode struct {
	CommonChaincode
}

func (t SideChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info(" Init ")
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t SideChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info(" Invoke ")
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Debug("fcn", fcn, "params", params)
	var responseBytes []byte
	switch fcn {
	case "get":
		var args = ArgsBuilder("get")
		var key = params[0]
		args.AppendArg(key)
		responseBytes = t.InvokeChaincode(mainCC, args.Get(), "").Payload
	}
	t.Logger.Debug("response", responseBytes)
	return shim.Success(responseBytes)
}

func main() {
	var cc = SideChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
