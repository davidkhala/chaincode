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
	testutil.AssertSame(t,response.Status,int32(200));
	fmt.Print("init ",response)

}
func TestStressChaincode_Invoke(t *testing.T) {
	var cc = new(StressChaincode)
	var stub = new(shim.ChaincodeStubInterface);
	var response = cc.Invoke(*stub);
	testutil.AssertSame(t,response.Status,int32(200));
	fmt.Print("invoke ",response)
}
func TestStart(t *testing.T) {
	shim.Start(new(StressChaincode))
}
