package main

import (
	"testing"
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"math/rand"
	"strconv"
	"fmt"
	"time"
)

var DavidExchange = ID{"david", "e"}
var LamMerchant = ID{"lam", "m"}
var StanleyConsumer = ID{"stanley,", "c"}

var purchaseTxID string

const (
	issueAmount = 100
	change      = 10
	cost        = 1
)

func TestPrepare(t *testing.T) {
	TestTradeChaincode_Init(t)
}
func randTxId(prefix string) string {
	return prefix + "|" + strconv.Itoa(rand.Int())
}
func TestTradeChaincode_InvokeCreate(t *testing.T) {
	var TxID = randTxId("create")
	var createArgs = [][]byte{[]byte(fcnWalletCreate), golang.ToJson(DavidExchange)}

	var response = mock.MockInvoke(TxID, createArgs)
	testutil.AssertSame(t, response.Status, int32(200));
	createArgs = [][]byte{[]byte(fcnWalletCreate), golang.ToJson(LamMerchant)}
	TxID = randTxId("create")
	response = mock.MockInvoke(TxID, createArgs)
	testutil.AssertSame(t, response.Status, int32(200));
	createArgs = [][]byte{[]byte(fcnWalletCreate), golang.ToJson(StanleyConsumer)}
	TxID = randTxId("create")
	response = mock.MockInvoke(TxID, createArgs)
}

func TestTradeChaincode_InvokeIssue(t *testing.T) {
	var TxID = randTxId("issue")
	var tx = CommonTransaction{ID{}, DavidExchange, issueAmount, tt_new_eToken_issue, 0}
	var invokeArgs = [][]byte{
		[]byte(tt_new_eToken_issue), //fcn
		golang.ToJson(DavidExchange),
		golang.ToJson(tx),
	}

	var response = mock.MockInvoke(TxID, invokeArgs)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestTradeChaincode_InvokeBalance(t *testing.T) {
	var TxID = randTxId("balance")
	var invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(DavidExchange),
	}

	var response = mock.MockInvoke(TxID, invokeArgs)
	var balanceResp BalanceResponse
	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertSame(t, balanceResp.Regular, int64(issueAmount))
}
func TestTradeChaincode_InvokeFiat(t *testing.T) {
	var TxID = randTxId("fiat")

	var tx = CommonTransaction{DavidExchange, StanleyConsumer, int64(change), tt_fiat_eToken_exchange, 0}
	var invokeArgs = [][]byte{
		[]byte(tt_fiat_eToken_exchange),
		golang.ToJson(DavidExchange),
		golang.ToJson(tx),
	}

	var response = mock.MockInvoke(TxID, invokeArgs)

	TxID = randTxId("balance")
	invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(DavidExchange),
	}
	response = mock.MockInvoke(TxID, invokeArgs)
	var balanceResp BalanceResponse
	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertNotEquals(t, balanceResp.Regular, int64(issueAmount))
	testutil.AssertEquals(t, balanceResp.Regular, int64(issueAmount-change))
}
func TestTradeChaincode_InvokeTimeStamp(t *testing.T) {
	var time1 = mock.TxTimestamp
	time.Sleep(1 * time.Second)
	var TxID = randTxId("timeStamp")
	var invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(DavidExchange),
	}
	mock.MockInvoke(TxID, invokeArgs)
	var time2 = mock.TxTimestamp
	fmt.Println("time1", time1, "time2", time2)
	testutil.AssertEquals(t, time2.Seconds > time1.Seconds, true)
}
func TestTradeChaincode_InvokePurchase(t *testing.T) {
	var tx = PurchaseTransaction{
		CommonTransaction{
			StanleyConsumer,
			LamMerchant,
			cost,
			tt_consumer_purchase,
			0,
		},
		"item1",
		1,
		"astri",
		"",
	}
	var TxID = randTxId("purchase")
	var balanceResp BalanceResponse
	var invokeArgs = [][]byte{
		[]byte(tt_consumer_purchase),
		golang.ToJson(StanleyConsumer),
		golang.ToJson(tx),
	}
	var response = mock.MockInvoke(TxID, invokeArgs)
	purchaseTxID = string(response.Payload) //IMPORTANT to use in accept purchase
	TxID = randTxId("consumer after purchase")
	invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(StanleyConsumer),
	}
	response = mock.MockInvoke(TxID, invokeArgs)
	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertEquals(t, balanceResp.Regular, int64(change-cost))
	TxID = randTxId("merchant after purchase")
	invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(LamMerchant),
	}
	response = mock.MockInvoke(TxID, invokeArgs)

	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertEquals(t, balanceResp.Escrow, int64(cost))
	testutil.AssertEquals(t, balanceResp.Regular, int64(0))
}
func TestTradeChaincode_InvokeAccept(t *testing.T) {
	var tx = PurchaseArbitrationTransaction{
		CommonTransaction{
			ID{},
			ID{},
			0,
			"",
			0,
		},
		true,
		purchaseTxID,
	}
	var invokeArgs = [][]byte{
		[]byte(tt_merchant_accept_purchase),
		golang.ToJson(LamMerchant),
		golang.ToJson(tx),
	}
	var TxID = randTxId("accept purchase ")

	var response = mock.MockInvoke(TxID, invokeArgs)

	TxID = randTxId("merchant after accept")
	invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(LamMerchant),
	}
	response = mock.MockInvoke(TxID, invokeArgs)
	var balanceResp BalanceResponse
	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertEquals(t, balanceResp.Regular, int64(cost))
}
func TestTradeChaincode_InvokePurchaseReject(t *testing.T) {
	var tx = PurchaseTransaction{
		CommonTransaction{
			StanleyConsumer,
			LamMerchant,
			cost,
			tt_consumer_purchase,
			0,
		},
		"item1",
		1,
		"astri",
		"",
	}
	var balanceResp BalanceResponse
	var TxID = randTxId("purchase")

	var invokeArgs = [][]byte{
		[]byte(tt_consumer_purchase),
		golang.ToJson(StanleyConsumer),
		golang.ToJson(tx),
	}
	var response = mock.MockInvoke(TxID, invokeArgs)
	purchaseTxID = string(response.Payload) //IMPORTANT to use in accept purchase

	var txReject = PurchaseArbitrationTransaction{
		CommonTransaction{
			ID{},
			ID{},
			0,
			"",
			0,
		},
		true,
		purchaseTxID,
	}
	invokeArgs = [][]byte{
		[]byte(tt_merchant_reject_purchase),
		golang.ToJson(LamMerchant),
		golang.ToJson(txReject),
	}
	TxID = randTxId("reject purchase ")

	mock.MockInvoke(TxID, invokeArgs)

	////
	TxID = randTxId("consumer after purchase reject")
	invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(StanleyConsumer),
	}
	response = mock.MockInvoke(TxID, invokeArgs)
	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertEquals(t, balanceResp.Regular, int64(change-cost))
	TxID = randTxId("merchant after purchase reject")
	invokeArgs = [][]byte{
		[]byte(fcnWalletBalance),
		golang.ToJson(LamMerchant),
	}
	response = mock.MockInvoke(TxID, invokeArgs)

	golang.FromJson(response.Payload, &balanceResp)
	testutil.AssertEquals(t, balanceResp.Regular, int64(cost))
}
