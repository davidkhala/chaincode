package main

import (
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"errors"
)

const (
	fcnWalletCreate  = "walletCreate"
	fcnWalletBalance = "walletBalance"
	fcnHistory       = "walletHistory"
	fcnListPurchase  = "listPurchase"

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

type PurchaseTransaction struct {
	CommonTransaction
	MerchandiseCode             string
	MerchandiseAmount           int64
	ConsumerDeliveryInstruction string
	Status                      string
}

func (tx PurchaseTransaction) isValid() {
	if tx.MerchandiseCode == "" {
		golang.PanicString("invalid PurchaseTransaction: empty MerchandiseCode")
	}
	if tx.MerchandiseAmount < 0 {
		golang.PanicString("invalid PurchaseTransaction: MerchandiseAmount<0")
	}
}

func (tx *PurchaseTransaction) Accept() (golang.Modifier) {
	return func(interface{}) {
		if tx.Status != StatusPending {
			golang.PanicString("Before accept purchase, invalid current status:" + tx.Status)
		}
		tx.Status = StatusAccepted
	}
}
func (tx *PurchaseTransaction) Reject() (golang.Modifier) {
	return func(interface{}) {
		if tx.Status != StatusPending {
			golang.PanicString("Before reject purchase, invalid current status:" + tx.Status)
		}
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
	if amount < 0 {
		golang.PanicString("invalid wallet value modification: amount<0")
	}
	return func(interface{}) {
		value.Balance += amount
		value.RecordID = recordID
	}
}
func (value *WalletValue) Lose(amount int64, recordID string, who string) (golang.Modifier) {
	if amount < 0 {
		golang.PanicString("invalid wallet value modification: amount<0")
	}
	return func(interface{}) {
		if value.Balance-amount < 0 {
			golang.PanicString(who + " has not enough Balance to pay " + golang.ToString(amount) + ", only have [" + golang.ToString(value.Balance) + "]")
		}
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

type HistoryPurchase struct {
	History map[string]PurchaseTransaction
}

type HistoryResponse struct {
	ID             ID
	RegularHistory []CommonTransaction
	EscrowHistory  []CommonTransaction
}
type Filter struct {
	Start  int64
	End    int64
	Status string
}
