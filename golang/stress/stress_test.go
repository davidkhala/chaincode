package main

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"fmt"
)

func TestStressChaincode_Init(t *testing.T) {
	var cc = new(StressChaincode)
	var stub = new(shim.ChaincodeStubInterface);
	var response peer.Response
	response = cc.Init(*stub)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestStressChaincode_Invoke(t *testing.T) {
	var cc = new(StressChaincode)
	var stub = new(shim.ChaincodeStubInterface);
	var response = cc.Invoke(*stub);
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
//	when error status is 500
}
func TestStart(t *testing.T) {
	shim.Start(new(StressChaincode))
}
