package main

import (
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"errors"
)

type ID struct {
	Name   string
	Prefix string
}

type WalletID struct {
	Name string
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
func (id ID) getWalletID() WalletID {
	var walletPrefix = "wallet_"

	return WalletID{walletPrefix + id.getLoginID()}
}
func (wallet WalletID) getEscrow() string {
	return wallet.Name + "_e"
}
func (wallet WalletID) getRegular() string {
	return wallet.Name + "_r"
}

