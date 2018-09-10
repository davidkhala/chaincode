package main

import (
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"strconv"
)

var cc = new(PrivateDataCC)
var mock = shim.NewMockStub(name, cc)

func TestStressChaincode_Init(t *testing.T) {

	var initArgs [][]byte
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
	var num, _ = strconv.Atoi(response.Message)
	testutil.AssertSame(t, num, 0);
}
func TestStressChaincode_Invoke(t *testing.T) {
	var args [][]byte

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
	var value = string(response.Payload)
	var num, _ = strconv.Atoi(value);
	testutil.AssertSame(t, num, 1)
	//	when error status is 500
}
func TestStressChaincode_Invoke1(t *testing.T) {
	var args [][]byte

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
	var value = string(response.Payload)
	var num, _ = strconv.Atoi(value);
	testutil.AssertSame(t, num, 2)
	//	when error status is 500
}
