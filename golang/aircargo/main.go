package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "airCargo"
)

type AirCargoChaincode struct {
	CommonChaincode
}

func (t AirCargoChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Init")
	return shim.Success(nil)
}

func (t AirCargoChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Invoke")
	var fcn, params = t.GetFunctionAndArgs()

	switch fcn {
	case "create":
		var req createAirCargo
		FromJson(params[0], &req)
		t.PutStateObj(req.AirCargoID, req)
		break;
	case "handler":
		break;
	case "transfer":
		break;
	default:
		PanicString("unknown fcn")
	}
	return shim.Success(nil)
}

func main() {
	var cc = AirCargoChaincode{}

	cc.SetLogger(name)
	ChaincodeStart(cc)

}
