package main

import (
	"testing"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/davidkhala/chaincode/golang/trade/golang"
)

var cc = &TradeChaincode{golang.CommonChaincode{Mock: true, Debug: true}}
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

	var set = OrgMap{
		ConsumerType: Org{Name: "Consumer", MSPID: "ConsumerMSP"},
		MerchantType: Org{Name: "Merchant", MSPID: "MerchantMSP"},
		ExchangeType: Org{Name: "Exchange", MSPID: "ExchangeMSP"},
	}
	var mspsBytes = golang.ToJson(&set)
	initArgs = append(initArgs, mspsBytes)
	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
