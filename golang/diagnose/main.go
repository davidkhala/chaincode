package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
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

func (t diagnoseChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info(" Init ")
	// GetStatus in Init will timeout request
	return shim.Success(nil)

}

func (t diagnoseChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	fcn, params := stub.GetFunctionAndParameters()
	t.Logger.Info("Invoke", fcn, params)
	var response []byte
	switch fcn {
	case "worldStates":
		var states = t.WorldStates("")
		response = ToJson(states)
	case "whoami":
		var cid = NewClientIdentity(stub)
		response = ToJson(cid)
	case "get":
		var key = params[0]
		response = t.GetState(key)
	case "put":
		var key = params[0]
		var value = params[1]
		t.PutState(key, []byte(value))
	}
	return shim.Success(response)
}

func main() {
	var cc = diagnoseChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
