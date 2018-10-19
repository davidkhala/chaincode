package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"time"
)

const (
	name   = "sideChaincode"
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
	case "put":
		var key1 = UnixMilliSecond(t.GetTxTime()).String()
		{
			var args = ArgsBuilder("put")
			var value = key1
			args.AppendArg(key1)
			args.AppendArg(value)
			t.InvokeChaincode(mainCC, args.Get(), "")
		}
		var key2 = UnixMilliSecond(time.Now()).String() + "1"
		{
			var args = ArgsBuilder("put")
			var value = key2
			args.AppendArg(key2)
			args.AppendArg(value)
			t.InvokeChaincode(mainCC, args.Get(), "")
		}
		responseBytes = []byte(key1 + "|" + key2)
	case "get":
		var args = ArgsBuilder("get")
		var key = params[0]
		args.AppendArg(key)
		responseBytes = t.InvokeChaincode(mainCC, args.Get(), "").Payload
	}
	return shim.Success(responseBytes)
}

func main() {
	var cc = SideChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
