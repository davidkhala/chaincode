package main

import (
	"github.com/davidkhala/fabric-common-chaincode/golang"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	tt_new_eToken_issue         = "tt_new_eToken_issuance"
	tt_fiat_eToken_exchange     = "tt_fiat_eToken_exchange"
	tt_consumer_purchase        = "tt_consumer_purchase"
	tt_merchant_reject_purchase = "tt_merchant_reject_purchase"
	tt_merchant_accept_purchase = "tt_merchant_accept_purchase"
	tt                          = "tt_unspecified"
)

type ID struct {
	Name   string
	Prefix string
}

type Wallet struct {
	ID string
}
type WalletValue struct {
	RecordID string
	Balance int64
}

func BuildID(data []byte) ID {
	var id ID
	golang.FromJson(data, &id);
	if id.Prefix != "c" && id.Prefix != "m" && id.Prefix != "e" {
		golang.PanicError(errors.New("invalide ID prefix " + id.Prefix))
	}
	return id
}
func (id ID) getLoginID() string {
	return id.Prefix + id.Name
}
func (id ID) getWallet(suffix string) Wallet {
	var walletPrefix = "wallet_"

	switch suffix {
	case "":
	case "_e": //for Escrow
	case "_r": //for Regular
	default:
		golang.PanicError(errors.New("invalid wallet suffix:" + suffix));
	}

	return Wallet{walletPrefix + id.getLoginID() + suffix}
}

func (wallet Wallet) GetHistory(ccAPI shim.ChaincodeStubInterface, ) {
	var key = wallet.ID;
	golang.HistoryToArray(golang.GetHistoryForKey(ccAPI,key))//TODO
}
