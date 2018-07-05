package main

import (
	"testing"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/davidkhala/chaincode/golang/trade/golang"
)

var cc = new(TradeChaincode)
var mock = shim.NewMockStub(name, cc)

func TestIsSubset(t *testing.T) {
	var first = []string{"1", "2"}
	var second = []string{"1", "2", "3"}
	testutil.AssertEquals(t, IsSubset(first, second), true)
}
func TestTradeChaincode_Init(t *testing.T) {
	var initArgs [][]byte
	initArgs = append(initArgs, []byte("Initfcn")) //fcn
	var TxID = "ob"

	var set = golang.StringList{}
	set.Put("ConsumerMSP")
	set.Put("MerchantMSP")
	set.Put("ExchangerMSP")
	var mspsBytes = golang.ToJson(&set)
	fmt.Println("mspsBytes", string(mspsBytes))
	initArgs = append(initArgs, mspsBytes)
	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestTradeChaincode_Invoke(t *testing.T) {
	var TxID = "oa"
	var invokeArgs [][]byte
	invokeArgs = append(invokeArgs, []byte("Invokefcn")) //fcn
	var response = mock.MockInvoke(TxID, invokeArgs)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
