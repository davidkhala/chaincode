package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/goutils"
)

const (
	FcnCreateToken  = "createToken"
	FcnGetToken     = "getToken"
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
	TokenCreateRequest
	Issuer       string // uses MSP ID in ecosystem
	Manager      string // uses MSP ID in ecosystem
	OwnerType    OwnerType
	IssuerClient cid.ClientIdentity
	TransferTime goutils.TimeLong
	Client       cid.ClientIdentity // latest Operator Client
}

type TokenCreateRequest struct {
	Owner    string
	MintTime goutils.TimeLong
	Content  []byte
}

func (t TokenCreateRequest) Build() TokenData {
	return TokenData{
		TokenCreateRequest: t,
	}
}

type TokenTransferRequest struct {
	Owner        string
	OwnerType    OwnerType
	TransferTime goutils.TimeLong
}

func (data *TokenData) Apply(request TokenTransferRequest) *TokenData {

	data.Owner = request.Owner
	data.OwnerType = request.OwnerType
	if data.TransferTime < request.TransferTime {
		data.TransferTime = request.TransferTime
	} else {
		panic("invalid TransferTime")
	}

	return data
}
