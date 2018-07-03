package main

import (
	"testing"
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
var cc = new(AdminChaincode)
var mock = shim.NewMockStub(name, cc)

func TestAdminChaincode_Init(t *testing.T) {
	var initArgs [][]byte
	initArgs = append(initArgs, []byte("Initfcn")) //fcn
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestAdminChaincode_Invoke(t *testing.T) {
	var args [][]byte
	args = append(args, []byte("Invokefcn")) //fcn

	var TxID = "oa"
	// stores a channel ID of the proposal
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}