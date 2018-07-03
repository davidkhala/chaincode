package main

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"fmt"
)

var cc = new(StressChaincode)
var mock = shim.NewMockStub(name, cc)

func TestStressChaincode_Init(t *testing.T) {
	var initArgs [][]byte
	initArgs = append(initArgs, []byte("Initfcn")) //fcn
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestStressChaincode_Invoke(t *testing.T) {

	var args [][]byte
	args = append(args, []byte("Invokefcn")) //fcn

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
	//	when error status is 500
}
func TestStart(t *testing.T) {
	shim.Start(new(StressChaincode))
}
