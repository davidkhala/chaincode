package main

import (
	"testing"
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/testutil"
)

var DavidExchange = ID{"david", "e"}
var LamMerchant = ID{"lam","m"}

const (
	issueAmount = 100
)

func TestPrepare(t *testing.T) {
	TestTradeChaincode_Init(t)
}

func TestTradeChaincode_InvokeCreate(t *testing.T) {
	var TxID = "create1"
	var createDavidArgs = [][]byte{[]byte(walletCreate),golang.ToJson(DavidExchange)}

	var response = mock.MockInvoke(TxID, createDavidArgs)
	testutil.AssertSame(t, response.Status, int32(200));
	var createLamArgs = [][]byte{[]byte(walletCreate),golang.ToJson(LamMerchant)}
	TxID = "create2"
	response = mock.MockInvoke(TxID, createLamArgs)
	testutil.AssertSame(t, response.Status, int32(200));
}

func TestTradeChaincode_InvokeIssue(t *testing.T) {
	var TxID = "issue"
	var tx = CommonTransaction{ID{}, DavidExchange, issueAmount, tt_new_eToken_issue, 0}
	var invokeArgs  = [][]byte{
		[]byte(tt_new_eToken_issue),//fcn
		golang.ToJson(DavidExchange),
		golang.ToJson(tx),
	}

	var response = mock.MockInvoke(TxID, invokeArgs)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestTradeChaincode_InvokeBalance(t *testing.T) {
	var TxID = "balance"
	var invokeArgs  = [][]byte{
		[]byte(walletBalance),
		golang.ToJson(DavidExchange),
	}

	var response = mock.MockInvoke(TxID, invokeArgs)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
	testutil.AssertSame(t, golang.ToInt(response.Payload), issueAmount)
}
