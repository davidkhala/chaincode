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
func (t SideChaincode) newTokenGetter() string {
	var key1 = RandString(4, "abcd")
	var args = ArgsBuilder("get")
	args.AppendArg(key1)
	var responseBytes = t.InvokeChaincode(mainCC, args.Get(), "").Payload
	if responseBytes == nil {
		return key1
	}
	return t.newTokenGetter()

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

		var key1 = t.newTokenGetter()
		{
			var args = ArgsBuilder("put")
			var value = UnixMilliSecond(time.Now()).String()
			args.AppendArg(key1)
			args.AppendArg(value)
			t.InvokeChaincode(mainCC, args.Get(), "")
		}

		responseBytes = []byte(key1)
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
