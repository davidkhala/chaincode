package main

import (
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"errors"
)

const (
	fcnWalletCreate  = "create"
	fcnWalletBalance = "balance"
	fcnHistory       = "history"

	tt_new_eToken_issue         = "tt_new_eToken_issuance"
	tt_fiat_eToken_exchange     = "tt_fiat_eToken_exchange"
	tt_consumer_purchase        = "tt_consumer_purchase"
	tt_merchant_reject_purchase = "tt_merchant_reject_purchase"
	tt_merchant_accept_purchase = "tt_merchant_accept_purchase"
	tt                          = "tt_unspecified"

	ConsumerMSP  = "ConsumerMSP"
	MerchantMSP  = "MerchantMSP"
	ExchangerMSP = "ExchangeMSP"

	ConsumerType   = "c"
	MerchantType   = "m"
	ExchangerType  = "e"
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
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
	MerchandiseCode             string
	MerchandiseAmount           int64
	ConsumerDeliveryInstruction string
	Status                      string
}

func (tx *PurchaseTransaction) Accept() (golang.Modifier) {
	if tx.Status != StatusPending {
		golang.PanicString("Before accept purchase, invalid current status:" + tx.Status)
	}
	return func(interface{}) {
		tx.Status = StatusAccepted
	}
}
func (tx *PurchaseTransaction) Reject() (golang.Modifier) {
	if tx.Status != StatusPending {
		golang.PanicString("Before reject purchase, invalid current status:" + tx.Status)
	}
	return func(interface{}) {
		tx.Status = StatusRejected
	}
}

type PurchaseArbitrationTransaction struct {
	CommonTransaction
	Accept       bool
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
