package main

import (
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	walletCreate                = "create"
	walletBalance               = "balance"
	tt_new_eToken_issue         = "tt_new_eToken_issuance"
	tt_fiat_eToken_exchange     = "tt_fiat_eToken_exchange"
	tt_consumer_purchase        = "tt_consumer_purchase"
	tt_merchant_reject_purchase = "tt_merchant_reject_purchase"
	tt_merchant_accept_purchase = "tt_merchant_accept_purchase"
	tt                          = "tt_unspecified"
)

type CommonTransaction struct {
	From      wallet
	To        wallet
	Amount    int64
	Type      string
	TimeStamp int64
}

type ID struct {
	Name   string
	Prefix string
}

type wallet struct {
	ID string
}
type WalletValue struct {
	RecordID string
	Balance  int64
}

func (id ID) getLoginID() string {
	return id.Prefix + id.Name
}
func (id ID) getWallet(suffix string) wallet {
	var walletPrefix = "wallet_"

	if id.Prefix != "c" && id.Prefix != "m" && id.Prefix != "e" {
		golang.PanicError(errors.New("invalide ID prefix " + id.Prefix))
	}
	switch suffix {
	case "":
	case "_e": //for Escrow
	case "_r": //for Regular
	default:
		golang.PanicError(errors.New("invalid wallet suffix:" + suffix));
	}

	return wallet{walletPrefix + id.getLoginID() + suffix}
}

func (wallet wallet) GetHistory(ccAPI shim.ChaincodeStubInterface, ) {
	var key = wallet.ID;
	golang.ParseHistory(golang.GetHistoryForKey(ccAPI, key))
}
