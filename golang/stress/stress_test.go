package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

var cc = new(StressChaincode)
var mock = shimtest.NewMockStub(name, cc)

func TestStressChaincode_Init(t *testing.T) {
	var initArgs [][]byte
	initArgs = append(initArgs, []byte("Initfcn")) // fcn
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
}
func TestStressChaincode_Invoke(t *testing.T) {

	var args [][]byte
	args = append(args, []byte("Invokefcn")) // fcn

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
}
func TestStart(t *testing.T) {
	shim.Start(new(StressChaincode))
}
