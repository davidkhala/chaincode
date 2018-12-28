package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "diagnose"
)

type diagnoseChaincode struct {
	CommonChaincode
}

func (t diagnoseChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Init", fcn, params)
	t.printTransient()
	return shim.Success(nil)

}
func (t diagnoseChaincode) printTransient() {
	var transientMap = t.GetTransient()
	t.Logger.Debug("==[start]transientMap")
	for k, v := range transientMap {
		t.Logger.Debug(k, ":", string(v))
	}
	t.Logger.Debug("==[end]transientMap")
}

func (t diagnoseChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	fcn, params := stub.GetFunctionAndParameters()
	t.Logger.Info("Invoke", fcn, params)
	t.printTransient()
	var responseBytes []byte
	switch fcn {
	case "panic":
		PanicString("test panic")
	case "richQuery":
		var query = params[0]
		t.Logger.Info("Query string", query)
		var queryIter = t.GetQueryResult(query)
		var states = ParseStates(queryIter)
		responseBytes = ToJson(states)
	case "worldStates":
		var states = t.WorldStates("")
		responseBytes = ToJson(states)
	case "whoami":
		responseBytes = ToJson(cid.NewClientIdentity(stub))
	case "get":
		var key = params[0]
		responseBytes = t.GetState(key)
	case "put":
		var key = params[0]
		var value = params[1]
		t.PutState(key, []byte(value))
	case "delegate":
		type crossChaincode struct {
			ChaincodeName string
			Fcn           string
			Args          []string
			Channel       string
		}
		var paramInput crossChaincode
		FromJson([]byte(params[0]), &paramInput)
		var args = ArgsBuilder(paramInput.Fcn)
		for i, element := range paramInput.Args {
			args.AppendArg(element)
			t.Logger.Debug("delegated Arg", i, element)
		}
		var pb = t.InvokeChaincode(paramInput.ChaincodeName, args.Get(), paramInput.Channel)
		responseBytes = pb.Payload
	}
	return shim.Success(responseBytes)
}

func main() {
	var cc = diagnoseChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
