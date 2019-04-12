package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "playGround"
)

type PlayGroundChaincode struct {
	CommonChaincode
}

func (t PlayGroundChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Init")
	return shim.Success(nil)
}

func (t PlayGroundChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Invoke")
	return shim.Success(nil)
}

func main() {
	var cc = PlayGroundChaincode{}

	cc.SetLogger(name)
	ChaincodeStart(cc)

}
