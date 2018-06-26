package main

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/common/ledger/testutil"
)

func TestStressChaincode_Init(t *testing.T) {
	var cc = new(StressChaincode)
	var stub = new(shim.ChaincodeStubInterface);
	var response peer.Response
	response = cc.Init(*stub)
	testutil.AssertEquals(t,response.Status,200);
	fmt.Print(response);
}
func TestStressChaincode_Invoke(t *testing.T) {
	var cc = new(StressChaincode)
	var stub = new(shim.ChaincodeStubInterface);
	var response = cc.Invoke(*stub);
	fmt.Print(response);
}
func TestStart(t *testing.T) {
	shim.Start(new(StressChaincode))
}
