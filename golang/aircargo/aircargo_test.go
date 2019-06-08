package main

import (
	"fmt"
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

var cc = new(AirCargoChaincode)
var mock = shim.NewMockStub(name, cc)

func TestStressChaincode_Init(t *testing.T) {
	var initArgs = ArgsBuilder("Initfcn").Get()
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
}
func TestStressChaincode_Invoke(t *testing.T) {

	var args = ArgsBuilder("Invokefcn").Get()

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
}
