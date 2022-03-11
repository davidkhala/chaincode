package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/goutils"
)

const (
	FcnCreateToken  = "createToken"
	FcnGetToken     = "getToken"
	FcnRenewToken   = "renewToken"
	FcnTokenHistory = "tokenHistory"
	FcnDeleteToken  = "deleteToken"
	FcnMoveToken    = "moveToken"
)

type OwnerType byte

const (
	_ = iota
	OwnerTypeMember
	OwnerTypeClinic
	OwnerTypeNetwork
	OwnerTypeInsurance
)

func (t OwnerType) To() string {
	var enum = []string{"unknown", "member", "clinic", "network", "insurance"}
	return enum[t]
}

type TokenData struct {
	Owner        string
	Issuer       string // uses MSP ID in ecosystem
	Manager      string // uses MSP ID in ecosystem
	OwnerType    OwnerType
	IssuerClient cid.ClientIdentity
	ExpiryDate   goutils.TimeLong
	TransferDate goutils.TimeLong
	Client       cid.ClientIdentity // latest Operator Client
}

type TokenCreateRequest struct {
	Owner      string
	ExpiryDate goutils.TimeLong
}

func (t TokenCreateRequest) Build() TokenData {
	return TokenData{
		Owner:      t.Owner,
		ExpiryDate: t.ExpiryDate,
	}
}

type TokenTransferRequest struct {
	Owner    string
	MetaData []byte
}

func (t TokenTransferRequest) ApplyOn(data TokenData) TokenData {
	if t.Owner != "" {
		data.Owner = t.Owner
	}

	return data
}
