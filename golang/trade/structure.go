package main

import (
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"errors"
)

const (
	walletCreate  = "create"
	walletBalance = "balance"

	tt_new_eToken_issue         = "tt_new_eToken_issuance"
	tt_fiat_eToken_exchange     = "tt_fiat_eToken_exchange"
	tt_consumer_purchase        = "tt_consumer_purchase"
	tt_merchant_reject_purchase = "tt_merchant_reject_purchase"
	tt_merchant_accept_purchase = "tt_merchant_accept_purchase"
	tt                          = "tt_unspecified"

	ConsumerMSP  = "ConsumerMSP"
	MerchantMSP  = "MerchantMSP"
	ExchangerMSP = "ExchangerMSP"

	ConsumerType  = "c"
	MerchantType  = "m"
	ExchangerType = "e"
)

type CommonTransaction struct {
	From      ID
	To        ID
	Amount    int64
	Type      string
	TimeStamp int64
}
type HistoryTransactions struct {
	History []CommonTransaction
}
type PurchaseTransaction struct {
	CommonTransaction
	Merchandise                 string
	MerchandiseAmount           int64
	ConsumerDeliveryInstruction string
}
type PurchaseArbitrationTransaction struct {
	CommonTransaction
	Status       bool
	PurchaseTxID string
}

type ID struct {
	Name string
	Type string
}

type wallet struct {
	regularID string
	escrowID  string
}
type WalletValue struct {
	RecordID string
	Balance  int64
}

func (value *WalletValue) Add(amount int64, recordID string) (golang.Modifier) {
	return func(interface{}) {
		value.Balance += amount
		value.RecordID = recordID
	}
}
func (value *WalletValue) Lose(amount int64, recordID string) (golang.Modifier) {
	return func(interface{}) {
		value.Balance -= amount
		value.RecordID = recordID
	}
}

func (id ID) getLoginID() string {
	return id.Type + id.Name
}
func (id ID) getWallet() wallet {
	var walletPrefix = "wallet_"
	if id.Type != ConsumerType && id.Type != MerchantType && id.Type != ExchangerType {
		golang.PanicError(errors.New("invalid ID prefix " + id.Type))
	}
	if id.Type == MerchantType {
		return wallet{
			walletPrefix + id.getLoginID() + "_r",
			walletPrefix + id.getLoginID() + "_e",
		}
	} else {
		return wallet{walletPrefix + id.getLoginID(), ""}
	}

}
