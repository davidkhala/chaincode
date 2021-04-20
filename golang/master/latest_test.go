package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

var cc = new(PrivateDataCC)
var mock = shimtest.NewMockStub(name, cc)

func TestStressChaincode_Init(t *testing.T) {

	var initArgs [][]byte
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
}
func TestStressChaincode_Invoke(t *testing.T) {
	var args [][]byte

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
	//	when error status is 500
}
func TestStressChaincode_Invoke1(t *testing.T) {
	var args [][]byte

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
}
