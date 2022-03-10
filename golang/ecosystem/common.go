package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/goutils"
)

const (
	GlobalCCID      = "global"
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
	TokenType    TokenType
	IssuerClient cid.ClientIdentity
	ExpiryDate   goutils.TimeLong
	TransferDate goutils.TimeLong
	Client       cid.ClientIdentity // latest Operator Client
	MetaData     []byte
}
type TokenType byte

const (
	_ = iota
	TokenTypeVerify
	TokenTypePay
)

func (t TokenType) To() string {
	var enum = []string{"verify", "pay"}
	return enum[t]
}
func (TokenType) From(s string) TokenType {
	var typeMap = map[string]TokenType{"verify": TokenTypeVerify, "pay": TokenTypePay}
	return typeMap[s]
}

type TokenCreateRequest struct {
	Owner      string
	TokenType  TokenType
	ExpiryDate goutils.TimeLong
	MetaData   []byte
}

func (t TokenCreateRequest) Build() TokenData {
	return TokenData{
		Owner:      t.Owner,
		TokenType:  t.TokenType,
		ExpiryDate: t.ExpiryDate,
		MetaData:   t.MetaData,
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
	if t.MetaData != nil {
		data.MetaData = t.MetaData
	}

	return data
}
